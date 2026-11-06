package main

import (
	"errors"
	"fmt"
	"sort"
	"sync"
	"time"
)

// Integrador: junta Mutex, RWMutex, select, WaitGroup e semáforo.
// Um "serviço" que busca URLs em paralelo, com cache, limite de
// concorrência, timeout por pedido e estatísticas.

var ErrTimeout = errors.New("timeout")

type Estatisticas struct {
	mu                                  sync.Mutex
	Sucessos, Falhas, Timeouts, DaCache int
}

func (e *Estatisticas) Somar(campo *int) {
	e.mu.Lock()
	defer e.mu.Unlock()
	*campo++
}

type Servico struct {
	buscar  func(string) (string, error)
	timeout time.Duration
	sem     chan struct{}
	cacheMu sync.RWMutex
	cache   map[string]string
	Estat   Estatisticas
}

func NovoServico(buscar func(string) (string, error), limite int, timeout time.Duration) *Servico {
	return &Servico{
		buscar:  buscar,
		timeout: timeout,
		sem:     make(chan struct{}, limite),
		cache:   make(map[string]string),
	}
}

// comTimeout: o buffer 1 evita que a goroutine fique presa se o select
// já tiver desistido.
func (s *Servico) comTimeout(url string) (string, error) {
	type resp struct {
		c   string
		err error
	}
	ch := make(chan resp, 1)
	go func() {
		c, err := s.buscar(url)
		ch <- resp{c, err}
	}()
	select {
	case r := <-ch:
		return r.c, r.err
	case <-time.After(s.timeout):
		return "", ErrTimeout
	}
}

func (s *Servico) um(url string) (string, error) {
	s.cacheMu.RLock()
	c, ok := s.cache[url]
	s.cacheMu.RUnlock()
	if ok {
		s.Estat.Somar(&s.Estat.DaCache)
		return c, nil
	}

	s.sem <- struct{}{}
	defer func() { <-s.sem }()

	c, err := s.comTimeout(url)
	switch {
	case errors.Is(err, ErrTimeout):
		s.Estat.Somar(&s.Estat.Timeouts)
	case err != nil:
		s.Estat.Somar(&s.Estat.Falhas)
	default:
		s.Estat.Somar(&s.Estat.Sucessos)
		s.cacheMu.Lock()
		s.cache[url] = c
		s.cacheMu.Unlock()
	}
	return c, err
}

// Todos devolve conteúdo por URL e o erro de cada URL que falhou.
func (s *Servico) Todos(urls []string) (map[string]string, map[string]error) {
	ok := make(map[string]string)
	erros := make(map[string]error)
	var mu sync.Mutex
	var wg sync.WaitGroup

	for _, url := range urls {
		wg.Add(1)
		go func() {
			defer wg.Done()
			c, err := s.um(url)
			mu.Lock()
			defer mu.Unlock()
			if err != nil {
				erros[url] = err
			} else {
				ok[url] = c
			}
		}()
	}
	wg.Wait()
	return ok, erros
}

func main() {
	buscar := func(u string) (string, error) {
		switch u {
		case "lento":
			time.Sleep(300 * time.Millisecond)
		case "mau":
			return "", errors.New("erro de rede")
		default:
			time.Sleep(10 * time.Millisecond)
		}
		return "conteúdo de " + u, nil
	}
	s := NovoServico(buscar, 2, 60*time.Millisecond)
	urls := []string{"a", "b", "c", "mau", "lento"}

	for ronda := 1; ronda <= 2; ronda++ {
		ok, erros := s.Todos(urls)
		var chaves []string
		for u := range erros {
			chaves = append(chaves, u)
		}
		sort.Strings(chaves)
		fmt.Printf("ronda %d: ok=%d erros=%v\n", ronda, len(ok), chaves)
	}
	e := &s.Estat
	fmt.Printf("sucessos=%d falhas=%d timeouts=%d cache=%d\n",
		e.Sucessos, e.Falhas, e.Timeouts, e.DaCache)
}
