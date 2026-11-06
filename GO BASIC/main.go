package main

import "fmt"

func Gerar(n int) <-chan int {

	ch := make(chan int)

	go func() {
		defer close(ch)
		for i := 1; i <= n; i++ {
			ch <- 1
		}
	}()

	return ch

}

func Consumir(in <-chan int) (soma int) {

	for v := range in {
		soma += v
	}

	return soma
}

func main() {
	fmt.Println("soma 1..5:", Consumir(Gerar(5)))
	fmt.Println("soma 1..0:", Consumir(Gerar(0)))

	// Channel com buffer: os valores ficam guardados depois do close.
	ch := make(chan int, 2)
	ch <- 10
	ch <- 20
	close(ch)

	v, ok := <-ch
	fmt.Println(v, ok) // 10 true
	v, ok = <-ch
	fmt.Println(v, ok) // 20 true
	v, ok = <-ch
	fmt.Println(v, ok) // 0 false  (fechado e vazio)

	func() {
		defer func() { fmt.Println("recuperado:", recover()) }()
		ch <- 30
	}()

}
