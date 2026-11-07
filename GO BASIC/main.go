package main

import (
	"fmt"
)

func Gerar(n int) <-chan int {

	out := make(chan int)

	go func() {
		defer close(out)
		for i := 1; i <= n; i++ {
			out <- i
		}
	}()

	return out
}

func Mapear(in <-chan int, f func(int) int) <-chan int {

	out := make(chan int)

	go func() {
		defer close(out)
		for v := range in {
			out <- f(v)
		}
	}()

	return out
}

func Filtrar(in <-chan int, ok func(int) bool) <-chan int {
	out := make(chan int)

	go func() {
		// se nao fecharmos  chanel , da deadlock
		defer close(out)
		for v := range in {
			if ok(v) {
				out <- v
			}
		}

	}()

	return out
}

func Somar(in <-chan int) int {
	total := 0
	for v := range in {
		total += v
	}

	return total
}

func main() {

	dobros := Mapear(Gerar(10), func(x int) int { return x * 2 })
	mult3 := Filtrar(dobros, func(x int) bool { return x%3 == 0 })
	fmt.Println("soma:", Somar(mult3))

	// Encadear tudo numa linha (lê-se de dentro para fora):
	fmt.Println(Somar(Filtrar(Mapear(Gerar(5), func(x int) int { return x * x }),
		func(x int) bool { return x%2 == 1 }))) // 1+9+25 = 35
}
