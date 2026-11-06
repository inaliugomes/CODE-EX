package main

import "fmt"

// somaParcial soma um pedaço da slice e envia o total pelo channel.
// "out chan<- int" = só pode ENVIAR (o compilador impede receber).
func somaParcial(pedaco []int, out chan<- int) {
	total := 0
	for _, v := range pedaco {
		total += v
	}
	out <- total
}

// SomaEmPartes generaliza o SomarParalelo (ex. 1) para N partes.
func SomaEmPartes(nums []int, partes int) int {
	if len(nums) == 0 || partes < 1 {
		return 0
	}
	if partes > len(nums) {
		partes = len(nums)
	}
	out := make(chan int, partes)            // buffer: nenhuma goroutine fica presa
	tam := (len(nums) + partes - 1) / partes // divisão arredondada para cima

	lancadas := 0
	for ini := 0; ini < len(nums); ini += tam {
		fim := min(ini+tam, len(nums))
		go somaParcial(nums[ini:fim], out)
		lancadas++
	}

	total := 0
	for i := 0; i < lancadas; i++ { // recebe exatamente tantas quantas lançou
		total += <-out
	}
	return total
}

// pingPong mostra channels SEM buffer: cada envio espera por um receptor.
func pingPong() {
	ping := make(chan string)
	pong := make(chan string)

	go func() {
		for msg := range ping { // termina quando ping for fechado
			pong <- msg + " -> recebido"
		}
		close(pong)
	}()

	for _, m := range []string{"a", "b", "c"} {
		ping <- m
		fmt.Println(<-pong)
	}
	close(ping)
	_, ok := <-pong // pong fechado: devolve valor zero e ok == false
	fmt.Println("pong ainda aberto?", ok)
}

// esperarTerminar usa um channel só para SINALIZAR (struct{} ocupa 0 bytes).
func esperarTerminar() {
	feito := make(chan struct{})
	go func() {
		fmt.Println("a trabalhar...")
		close(feito) // fechar acorda quem estiver à espera
	}()
	<-feito
	fmt.Println("terminou")
}

func main() {
	nums := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	fmt.Println(SomaEmPartes(nums, 3), SomaEmPartes(nums, 1), SomaEmPartes(nums, 50), SomaEmPartes(nil, 4))
	pingPong()
	esperarTerminar()
}
