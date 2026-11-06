package main

import (
	"fmt"
	"sort"
	"sync"
)

// Merge (fan-in) junta vários channels num só.
// Uma goroutine por entrada copia para a saída; uma goroutine extra
// fecha a saída quando TODAS as cópias acabarem.
func Merge(entradas ...<-chan int) <-chan int {
	out := make(chan int)
	var wg sync.WaitGroup

	copiar := func(c <-chan int) {
		defer wg.Done()
		for v := range c {
			out <- v
		}
	}

	wg.Add(len(entradas))
	for _, c := range entradas {
		go copiar(c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func gerarDe(inicio, fim int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for i := inicio; i <= fim; i++ {
			out <- i
		}
	}()
	return out
}

func main() {
	var todos []int
	for v := range Merge(gerarDe(1, 3), gerarDe(10, 12), gerarDe(100, 101)) {
		todos = append(todos, v)
	}
	sort.Ints(todos) // a ordem de chegada varia; ordenamos para comparar
	fmt.Println(todos)
}
