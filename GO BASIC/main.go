package main

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

func consumir(in <-chan int) (soma int) {

	for v := range in {
		soma += v
	}

	return soma
}

func main() {

}
