package main

import (
	"fmt"
	"sync"
)

// QuadradosParalelo: cada goroutine escreve só na SUA posição i,
// por isso não há corrida de dados (índices diferentes).
func QuadradosParalelo(nums []int) []int {
	res := make([]int, len(nums))
	var wg sync.WaitGroup

	for i, v := range nums {
		wg.Add(1) // ANTES de lançar a goroutine, nunca dentro dela
		go func() {
			defer wg.Done()
			res[i] = v * v // Go >= 1.22: i e v são novos em cada volta
		}()
	}
	wg.Wait()
	return res
}

// Resultado guarda valor E erro por posição: assim não se perde a ordem.
type Resultado struct {
	Valor int
	Erro  error
}

// DividirTodos faz 100/d para cada d. Divisão por zero vira erro, não panic.
func DividirTodos(divisores []int) []Resultado {
	res := make([]Resultado, len(divisores))
	var wg sync.WaitGroup

	for i, d := range divisores {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if d == 0 {
				res[i] = Resultado{Erro: fmt.Errorf("divisor %d: divisão por zero", i)}
				return
			}
			res[i] = Resultado{Valor: 100 / d}
		}()
	}
	wg.Wait()
	return res
}

// trabalhador recebe o WaitGroup POR PONTEIRO. Passar por valor copia-o
// e o Wait ficaria à espera para sempre (deadlock).
func trabalhador(id int, wg *sync.WaitGroup, saida []string) {
	defer wg.Done()
	saida[id] = fmt.Sprintf("trabalhador %d feito", id)
}

func main() {
	fmt.Println(QuadradosParalelo([]int{1, 2, 3, 4, 5}))

	for _, r := range DividirTodos([]int{5, 0, 20, 3}) {
		fmt.Println(r.Valor, r.Erro)
	}

	var wg sync.WaitGroup
	saida := make([]string, 3)
	for i := range 3 { // range sobre int (Go >= 1.22)
		wg.Add(1)
		go trabalhador(i, &wg, saida)
	}
	wg.Wait()
	fmt.Println(saida)
}
