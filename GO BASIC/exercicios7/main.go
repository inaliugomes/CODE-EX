package main

import (
	"sync"
)

// =====================================================================
// Exercício 1 — goroutines + channel
// SomarParalelo divide nums em duas metades, soma cada metade numa
// goroutine diferente, e devolve a soma total.
// Cada goroutine envia a sua soma parcial por um channel; a função
// recebe as duas e junta-as.
// Casos-limite: slice vazia (nil), com 1 elemento, com tamanho ímpar.
// Dica: meio := len(nums) / 2; nums[:meio] e nums[meio:].
// =====================================================================
func somar(nums []int, ch chan int) {

	var total int

	for _, v := range nums {

		total = total + v
	}

	ch <- total

}
func SomarParalelo(nums []int) int {

	if len(nums) == 0 {

		return 0
	}

	var total int
	meio := len(nums) / 2
	ch := make(chan int, 2)
	primeiraMetade := nums[:meio]
	segundaMetade := nums[meio:]

	go somar(primeiraMetade, ch)
	go somar(segundaMetade, ch)

	total = total + <-ch
	total = total + <-ch

	return total
}

// =====================================================================
// Exercício 2 — produtor (close + range)
// Gerar devolve um channel só de leitura que emite 1, 2, ..., n e depois
// é FECHADO (quem lê com range precisa disso para terminar).
// Gerar(0) devolve um channel que fecha logo, sem valores.
// Dica: cria o channel, lança uma goroutine que envia e faz close, e
// devolve o channel logo (sem esperar pela goroutine).
// =====================================================================
func Gerar(n int) <-chan int {

	ch := make(chan int)

	go func() {

		for i := 1; i < n+1; i++ {

			ch <- i
		}

		close(ch)

	}()

	return ch
}

// =====================================================================
// Exercício 3 — sync.WaitGroup
// QuadradosParalelo devolve uma slice com o quadrado de cada número,
// na MESMA ordem, calculando cada quadrado numa goroutine própria.
// Cada goroutine escreve só na sua posição i do resultado, por isso não
// há corrida de dados. Espera por todas com WaitGroup (Add, Done, Wait).
// Cuidado: Add antes de lançar a goroutine, nunca dentro dela.
// =====================================================================

func quadrado(num int, numero int, nums []int, wg *sync.WaitGroup) {

	defer wg.Done()

	nums[num] = numero * numero
}

func QuadradosParalelo(nums []int) []int {

	var wg sync.WaitGroup

	sliceResultado := make([]int, len(nums))

	for i, v := range nums {

		wg.Add(1)
		go quadrado(i, v, sliceResultado, &wg)
	}
	wg.Wait()

	return sliceResultado
}

// =====================================================================
// Exercício 4 — pipeline (etapa que lê de um channel e escreve noutro)
// Dobrar lê valores de in, e envia o dobro de cada um para o channel que
// devolve. Quando in fechar, o channel devolvido também tem de fechar.
// Dica: tal como no Gerar, devolve o channel logo e trabalha numa
// goroutine; for v := range in { ... } termina sozinho com o close.
// =====================================================================

func Dobrar(in <-chan int) <-chan int {

	ch := make(chan int)

	go func() {
		for v := range in {

			ch <- v * 2
		}
		close(ch)

	}()

	return ch
}

// =====================================================================
// Exercício 5 — worker pool (channel com buffer)
// Trabalhadores lança n goroutines ("trabalhadores") que leem tarefas de
// um channel e devolvem tarefa*10 por outro channel de resultados.
// Devolve todos os resultados (em qualquer ordem).
//
// Passos:
//  1. jobs e resultados com buffer len(tarefas) (assim ninguém bloqueia).
//  2. lança n trabalhadores: for j := range jobs { resultados <- j*10 }
//  3. envia as tarefas para jobs e fecha jobs.
//  4. espera pelos trabalhadores (WaitGroup) e só depois fecha resultados.
//  5. junta tudo o que está em resultados numa slice.
//
// Pergunta para ti: porque é que fechar jobs é obrigatório?
// =====================================================================
func Trabalhadores(tarefas []int, n int) []int {


	
	jobs := make(chan int , len(tarefas))
	resultado := make(chan int , len(tarefas))

	var wg sync.WaitGroup

	for i :=0 ; i < n ; i++ {

		wg.Add(1)

		go func(){

			defer wg.Done()

			for j := range jobs {

				resultado <- j * 10 
			}
		}()
	}

	for _, t := range tarefas {

		jobs <- t 
	}
	close(jobs)

	wg.Wait()
	close(resultado)

	var final []int 

	for r := range resultado {
		final = append(final,r)
	}

	return final 
}

func main() {}
