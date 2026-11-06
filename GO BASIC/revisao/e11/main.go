package main

import (
	"errors"
	"fmt"
	"time"
)

var ErrTimeout = errors.New("timeout")

// PrimeiroAResponder (ex. 4): devolve o primeiro valor que chegar.
func PrimeiroAResponder(a, b <-chan string) string {
	select {
	case v := <-a:
		return v
	case v := <-b:
		return v
	}
}

// Corrida generaliza para N fontes + limite de tempo.
// O channel tem buffer len(fontes): as goroutines mais lentas conseguem
// enviar e terminar mesmo que ninguém leia (sem goroutine leak).
func Corrida(fontes []func() string, limite time.Duration) (string, error) {
	ch := make(chan string, len(fontes))
	for _, f := range fontes {
		go func() { ch <- f() }()
	}
	select {
	case v := <-ch:
		return v, nil
	case <-time.After(limite):
		return "", ErrTimeout
	}
}

func fonte(nome string, d time.Duration) func() string {
	return func() string {
		time.Sleep(d)
		return nome
	}
}

func main() {
	a := make(chan string, 1)
	b := make(chan string, 1)
	go func() { time.Sleep(50 * time.Millisecond); a <- "lento" }()
	go func() { time.Sleep(5 * time.Millisecond); b <- "rápido" }()
	fmt.Println(PrimeiroAResponder(a, b))

	fontes := []func() string{
		fonte("A", 30*time.Millisecond),
		fonte("B", 10*time.Millisecond),
		fonte("C", 20*time.Millisecond),
	}
	fmt.Println(Corrida(fontes, time.Second))
	fmt.Println(Corrida(fontes, 2*time.Millisecond))
}
