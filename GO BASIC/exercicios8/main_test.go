package main

import (
	"errors"
	"fmt"
	"sort"
	"sync"
	"testing"
	"time"
)

func comTimeout(t *testing.T, f func()) {
	t.Helper()
	feito := make(chan struct{})
	go func() {
		f()
		close(feito)
	}()
	select {
	case <-feito:
	case <-time.After(2 * time.Second):
		t.Fatal("timeout: a função ficou bloqueada")
	}
}

func TestContador(t *testing.T) {
	comTimeout(t, func() {
		var c Contador
		var wg sync.WaitGroup
		for i := 0; i < 100; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for j := 0; j < 100; j++ {
					c.Incrementar()
				}
			}()
		}
		wg.Wait()
		if got := c.Valor(); got != 10000 {
			t.Errorf("Valor() = %d; esperado 10000", got)
		}
	})
}

func TestCache(t *testing.T) {
	comTimeout(t, func() {
		c := NovaCache()
		if _, ok := c.Obter("x"); ok {
			t.Error("chave inexistente devolveu ok=true")
		}
		var wg sync.WaitGroup
		for i := 0; i < 50; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				c.Pôr(fmt.Sprint("k", i), i)
				c.Obter("k0")
			}()
		}
		wg.Wait()
		for i := 0; i < 50; i++ {
			if v, ok := c.Obter(fmt.Sprint("k", i)); !ok || v != i {
				t.Errorf("Obter(k%d) = %d, %v; esperado %d, true", i, v, ok, i)
			}
		}
	})
}

func TestEsperarOuTimeout(t *testing.T) {
	t.Run("a tempo", func(t *testing.T) {
		ch := make(chan int, 1)
		ch <- 42
		v, err := EsperarOuTimeout(ch, 500*time.Millisecond)
		if err != nil || v != 42 {
			t.Errorf("got %d, %v; esperado 42, nil", v, err)
		}
	})
	t.Run("timeout", func(t *testing.T) {
		ch := make(chan int)
		inicio := time.Now()
		_, err := EsperarOuTimeout(ch, 50*time.Millisecond)
		if !errors.Is(err, ErrTimeout) {
			t.Errorf("err = %v; esperado ErrTimeout", err)
		}
		if time.Since(inicio) > time.Second {
			t.Error("demorou demasiado a dar timeout")
		}
	})
}

func TestPrimeiroAResponder(t *testing.T) {
	comTimeout(t, func() {
		a := make(chan string, 1)
		b := make(chan string, 1)
		go func() { time.Sleep(300 * time.Millisecond); a <- "lento" }()
		go func() { time.Sleep(20 * time.Millisecond); b <- "rápido" }()
		if got := PrimeiroAResponder(a, b); got != "rápido" {
			t.Errorf("got %q; esperado \"rápido\"", got)
		}
	})
}

func TestBuscarTodos(t *testing.T) {
	comTimeout(t, func() {
		buscar := func(u string) (string, error) {
			time.Sleep(10 * time.Millisecond)
			if u == "mau1" || u == "mau2" {
				return "", errors.New("falhou")
			}
			return "conteúdo de " + u, nil
		}
		urls := []string{"a", "mau1", "b", "c", "mau2"}
		res, falhas := BuscarTodos(urls, buscar)
		if len(res) != 3 || res["a"] != "conteúdo de a" || res["c"] != "conteúdo de c" {
			t.Errorf("res = %v", res)
		}
		sort.Strings(falhas)
		if len(falhas) != 2 || falhas[0] != "mau1" || falhas[1] != "mau2" {
			t.Errorf("falhas = %v", falhas)
		}
	})
}
