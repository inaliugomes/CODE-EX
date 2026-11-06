package main

import (
	"errors"
	"fmt"
	"sort"
	"time"
)

var ErrTimeout = errors.New("timeout")

// 1) select com timeout.
func EsperarOuTimeout(ch <-chan int, d time.Duration) (int, error) {
	select {
	case v := <-ch:
		return v, nil
	case <-time.After(d):
		return 0, ErrTimeout
	}
}

// 2) default = não bloqueia (tenta e segue em frente).
func TentarReceber(ch <-chan int) (int, bool) {
	select {
	case v := <-ch:
		return v, true
	default:
		return 0, false
	}
}

func TentarEnviar(ch chan<- int, v int) bool {
	select {
	case ch <- v:
		return true
	default: // buffer cheio ou ninguém a ler
		return false
	}
}

// 3) channel nil NUNCA fica pronto: serve para "desligar" um case.
// Lê de a e b até os dois fecharem.
func Juntar(a, b <-chan int) []int {
	var res []int
	for a != nil || b != nil {
		select {
		case v, ok := <-a:
			if !ok {
				a = nil // desliga este case
				continue
			}
			res = append(res, v)
		case v, ok := <-b:
			if !ok {
				b = nil
				continue
			}
			res = append(res, v)
		}
	}
	return res
}

func emitir(vals ...int) <-chan int {
	ch := make(chan int)
	go func() {
		defer close(ch)
		for _, v := range vals {
			ch <- v
		}
	}()
	return ch
}

// 4) for + select + canal "done" para parar uma goroutine de fora.
func contarTicks(parar <-chan struct{}) int {
	t := time.NewTicker(10 * time.Millisecond)
	defer t.Stop()
	n := 0
	for {
		select {
		case <-t.C:
			n++
		case <-parar:
			return n
		}
	}
}

func main() {
	rapido := make(chan int, 1)
	rapido <- 7
	fmt.Println(EsperarOuTimeout(rapido, time.Second))
	fmt.Println(EsperarOuTimeout(make(chan int), 20*time.Millisecond))

	vazio := make(chan int, 1)
	fmt.Println(TentarReceber(vazio))
	fmt.Println(TentarEnviar(vazio, 1), TentarEnviar(vazio, 2)) // true false

	r := Juntar(emitir(1, 2, 3), emitir(10, 20))
	sort.Ints(r)
	fmt.Println(r)

	parar := make(chan struct{})
	res := make(chan int)
	go func() { res <- contarTicks(parar) }()
	time.Sleep(65 * time.Millisecond)
	close(parar)
	fmt.Println("ticks >= 4:", <-res >= 4)
}
