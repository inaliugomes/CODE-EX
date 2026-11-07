package main

import (
	"fmt"
	"sync"
)

func QuadradosParalelo(nums []int) []int {

	res := make([]int , len(nums))
	var wg sync.WaitGroup
	
	for i , v := range nums {

		wg.Add(1)
		go func() {
			defer wg.Done()
			res[i] = v * v 
		}()

		wg.Wait()
		
	}
	return res 
}

type Resultado struct { 

	Valor int 
	Erro error 
}

func DividirTodos(divisores []int) []Resultado { 

	res := make([]Resultado, len(divisores))

	var wg sync.WaitGroup

	for i, d := range divisores { 
		wg.Add(1)
		go func() { 
			defer wg.Done()
			if d == 0 { 
				res[i] = Resultado{Erro: fmt.Errorf("%w")}
				return
			}
			res[i] = Resultado{Valor: 100/d}
		}()
	}

	wg.Wait()
	return res 

}

func trabalhador(id int , wg *sync.WaitGroup, saida []string){
	defer wg.Done()
	saida[id] = fmt.Sprintf("trabalhador %d feito00",id)
}