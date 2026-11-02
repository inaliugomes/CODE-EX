package main


import(
	"sync"
)

func main()  {
	
}

func Trabalhadores(tarefas []int, n int) []int {
        jobs := make(chan int, len(tarefas))
        resultados := make(chan int, len(tarefas))

        // 1. lançar n trabalhadores (peças A e B)
        var wg sync.WaitGroup
        for i := 0; i < n; i++ {
                wg.Add(1) // antes do go
                go func() {
                        defer wg.Done()
                        for j := range jobs { // termina quando jobs fechar
                                resultados <- j * 10
                        }
                }()
        }

        // 2. enviar as tarefas e fechar jobs
        for _, t := range tarefas {
                jobs <- t
        }
        close(jobs)

        // 3. esperar pelos trabalhadores, só depois fechar resultados
        wg.Wait()
        close(resultados)

        // 4. juntar tudo numa slice
        var final []int
        for r := range resultados {
                final = append(final, r)
        }
        return final
}
