package main

import (
	"fmt"
	"sync"
)

// Cache genérica com RWMutex: muitos leitores OU um escritor.
type Cache[V any] struct {
	mu    sync.RWMutex
	dados map[string]V
}

func NovaCache[V any]() *Cache[V] {
	return &Cache[V]{dados: make(map[string]V)} // map nil não aceita escrita
}

func (c *Cache[V]) Obter(k string) (V, bool) {
	c.mu.RLock() // leitura partilhada
	defer c.mu.RUnlock()
	v, ok := c.dados[k]
	return v, ok
}

func (c *Cache[V]) Por(k string, v V) {
	c.mu.Lock() // escrita exclusiva
	defer c.mu.Unlock()
	c.dados[k] = v
}

func (c *Cache[V]) Apagar(k string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.dados, k)
}

func (c *Cache[V]) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.dados)
}

// ObterOuCalcular usa "verificação dupla": o caminho rápido só lê (RLock);
// se falhar, sobe para Lock e volta a verificar, porque entre o RUnlock e o
// Lock outra goroutine pode ter chegado primeiro.
func (c *Cache[V]) ObterOuCalcular(k string, calcular func() V) V {
	if v, ok := c.Obter(k); ok {
		return v
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if v, ok := c.dados[k]; ok { // segunda verificação
		return v
	}
	v := calcular()
	c.dados[k] = v
	return v
}

func main() {
	c := NovaCache[int]()
	c.Por("a", 1)
	v, ok := c.Obter("a")
	fmt.Println(v, ok)
	_, ok = c.Obter("zzz")
	fmt.Println(ok)

	calculos := 0 // só mexido dentro de calcular(), que corre sob Lock
	var wg sync.WaitGroup
	for range 50 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			c.ObterOuCalcular("pesado", func() int {
				calculos++
				return 42
			})
		}()
	}
	wg.Wait()
	fmt.Println("calculos:", calculos, "len:", c.Len()) // 1 vez só

	c.Apagar("a")
	fmt.Println("len:", c.Len())
}
