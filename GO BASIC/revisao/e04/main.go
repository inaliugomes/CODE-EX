package main

import "fmt"

// Pipeline: cada etapa lê de um channel e devolve outro.
// Quando a entrada fecha, a etapa fecha a sua saída (propaga o fim).

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

// Mapear aplica f a cada valor (etapa genérica).
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

// Filtrar deixa passar só os valores onde ok(v) é verdadeiro.
func Filtrar(in <-chan int, ok func(int) bool) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for v := range in {
			if ok(v) {
				out <- v
			}
		}
	}()
	return out
}

// Somar é a etapa final: consome tudo e devolve um valor.
func Somar(in <-chan int) int {
	total := 0
	for v := range in {
		total += v
	}
	return total
}

func main() {
	// 1..10 -> dobro -> só os múltiplos de 3 -> soma
	dobros := Mapear(Gerar(10), func(x int) int { return x * 2 })
	mult3 := Filtrar(dobros, func(x int) bool { return x%3 == 0 })
	fmt.Println("soma:", Somar(mult3)) // 6+12+18 = 36

	// Encadear tudo numa linha (lê-se de dentro para fora):
	fmt.Println(Somar(Filtrar(Mapear(Gerar(5), func(x int) int { return x * x }),
		func(x int) bool { return x%2 == 1 }))) // 1+9+25 = 35
}
