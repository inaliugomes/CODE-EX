package main

import (
	"errors"
	"fmt"
	"sort"
	"sync"
	"sync/atomic"
	"time"
)

type Buscar func(url string) (string, error)

// Versão 1 — Mutex: cada goroutine escreve no map/slice partilhados.
// O download (buscar) fica FORA do lock: só o acesso aos dados é protegido.
func BuscarMutex(urls []string, buscar Buscar) (map[string]string, []string) {
	res := make(map[string]string)
	var falhas []string
	var mu sync.Mutex
	var wg sync.WaitGroup

	for _, url := range urls {
		wg.Add(1)
		go func() {
			defer wg.Done()
			conteudo, err := buscar(url)
			mu.Lock()
			defer mu.Unlock()
			if err != nil {
				falhas = append(falhas, url)
				return
			}
			res[url] = conteudo
		}()
	}
	wg.Wait()
	return res, falhas
}

// Versão 2 — channel: as goroutines só ENVIAM; uma única goroutine (a
// função) escreve no map. Sem Mutex porque nada é partilhado.
func BuscarCanal(urls []string, buscar Buscar) (map[string]string, []string) {
	type resp struct {
		url, conteudo string
		err           error
	}
	ch := make(chan resp, len(urls))
	var wg sync.WaitGroup

	for _, url := range urls {
		wg.Add(1)
		go func() {
			defer wg.Done()
			c, err := buscar(url)
			ch <- resp{url, c, err}
		}()
	}
	go func() { wg.Wait(); close(ch) }()

	res := make(map[string]string)
	var falhas []string
	for r := range ch {
		if r.err != nil {
			falhas = append(falhas, r.url)
		} else {
			res[r.url] = r.conteudo
		}
	}
	return res, falhas
}

// Versão 3 — limite de concorrência com um "semáforo" (channel com buffer).
// No máximo `limite` downloads ao mesmo tempo.
func BuscarLimitado(urls []string, limite int, buscar Buscar) (map[string]string, []string) {
	sem := make(chan struct{}, limite)
	res := make(map[string]string)
	var falhas []string
	var mu sync.Mutex
	var wg sync.WaitGroup

	for _, url := range urls {
		wg.Add(1)
		go func() {
			defer wg.Done()
			sem <- struct{}{}        // pede um lugar (bloqueia se estiver cheio)
			defer func() { <-sem }() // devolve o lugar
			c, err := buscar(url)
			mu.Lock()
			defer mu.Unlock()
			if err != nil {
				falhas = append(falhas, url)
			} else {
				res[url] = c
			}
		}()
	}
	wg.Wait()
	return res, falhas
}

func main() {
	buscar := func(u string) (string, error) {
		time.Sleep(10 * time.Millisecond)
		if u == "mau1" || u == "mau2" {
			return "", errors.New("falhou")
		}
		return "conteúdo de " + u, nil
	}
	urls := []string{"a", "mau1", "b", "c", "mau2"}

	for nome, f := range map[string]func([]string, Buscar) (map[string]string, []string){
		"mutex": BuscarMutex, "canal": BuscarCanal,
	} {
		res, falhas := f(urls, buscar)
		sort.Strings(falhas)
		fmt.Println(nome, len(res), falhas)
	}

	// Prova do limite: medimos quantos buscar() correm ao mesmo tempo.
	var ativos, maximo atomic.Int32
	medido := func(u string) (string, error) {
		n := ativos.Add(1)
		for {
			m := maximo.Load()
			if n <= m || maximo.CompareAndSwap(m, n) {
				break
			}
		}
		time.Sleep(10 * time.Millisecond)
		ativos.Add(-1)
		return u, nil
	}
	muitos := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"}
	BuscarLimitado(muitos, 3, medido)
	fmt.Println("máximo simultâneo <= 3:", maximo.Load() <= 3)
}
