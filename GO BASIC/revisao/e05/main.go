package main

import (
	"errors"
	"fmt"
	"sync"
)

type Trabalho struct {
	Indice int
	Valor  int
}

type Resultado struct {
	Indice int
	Valor  int
	Erro   error
}

// Processar corre f sobre cada tarefa com n trabalhadores e devolve os
// resultados NA ORDEM ORIGINAL (usa o Indice para recolocar).
func Processar(tarefas []int, n int, f func(int) (int, error)) []Resultado {
	jobs := make(chan Trabalho, len(tarefas))
	saida := make(chan Resultado, len(tarefas))

	var wg sync.WaitGroup
	for range n {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for t := range jobs { // acaba quando jobs fecha e esvazia
				v, err := f(t.Valor)
				saida <- Resultado{Indice: t.Indice, Valor: v, Erro: err}
			}
		}()
	}

	for i, v := range tarefas {
		jobs <- Trabalho{Indice: i, Valor: v}
	}
	close(jobs) // OBRIGATÓRIO: sem isto os range dos trabalhadores nunca acabam

	wg.Wait()
	close(saida) // só depois de TODOS os trabalhadores terminarem

	res := make([]Resultado, len(tarefas))
	for r := range saida {
		res[r.Indice] = r
	}
	return res
}

func main() {
	f := func(x int) (int, error) {
		if x < 0 {
			return 0, errors.New("negativo")
		}
		return x * 10, nil
	}
	for _, r := range Processar([]int{1, 2, -3, 4, 5, 6}, 3, f) {
		fmt.Printf("#%d valor=%d erro=%v\n", r.Indice, r.Valor, r.Erro)
	}
}
