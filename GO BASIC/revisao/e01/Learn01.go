package main

import "sync"


func BuscarTodos(urls []string, buscar func(string) (string, error)) (map[string]string, []string) {

	var wg sync.WaitGroup
	var mu sync.Mutex


	final := make(map[string]string)

	var erro []string 

	for _,url := range urls {

		wg.Add(1)

		go func(){

			defer wg.Done()

			resulta , err := buscar(url)
			mu.Lock()
			defer mu.Unlock()

			if err == nil { 

				final[url] = resulta
			}else { 
				erro = append(erro, url)
			}
			
		}()
	}
	wg.Wait()

	return	final , erro 
}