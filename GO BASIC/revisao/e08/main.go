package main

import (
	"errors"
	"fmt"
	"sync"
)

// ---- Contador: a regra de ouro é Lock + defer Unlock na linha seguinte.
type Contador struct {
	mu    sync.Mutex
	valor int
}

func (c *Contador) Incrementar() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.valor++ // ler-somar-gravar: sem lock, duas goroutines perdem somas
}

func (c *Contador) Valor() int {
	c.mu.Lock() // também a LEITURA precisa de lock
	defer c.mu.Unlock()
	return c.valor
}

// ---- Duas contas: bloquear DOIS mutexes pode dar deadlock.
// A faz Lock(A) e espera B; B faz Lock(B) e espera A -> ficam presas.
// Solução: bloquear SEMPRE pela mesma ordem (aqui, pelo ID).
var ErrSaldo = errors.New("saldo insuficiente")

type Conta struct {
	mu    sync.Mutex
	ID    int
	Saldo int
}

func Transferir(de, para *Conta, valor int) error {
	if de == para {
		return errors.New("mesma conta")
	}
	primeiro, segundo := de, para
	if para.ID < de.ID {
		primeiro, segundo = para, de
	}
	primeiro.mu.Lock()
	defer primeiro.mu.Unlock()
	segundo.mu.Lock()
	defer segundo.mu.Unlock()

	if de.Saldo < valor {
		return ErrSaldo
	}
	de.Saldo -= valor
	para.Saldo += valor
	return nil
}

func main() {
	var c Contador
	var wg sync.WaitGroup
	for range 1000 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range 100 {
				c.Incrementar()
			}
		}()
	}
	wg.Wait()
	fmt.Println("contador:", c.Valor()) // 100000

	a := &Conta{ID: 1, Saldo: 1000}
	b := &Conta{ID: 2, Saldo: 1000}
	for range 100 {
		wg.Add(2)
		go func() { defer wg.Done(); Transferir(a, b, 10) }() // A -> B
		go func() { defer wg.Done(); Transferir(b, a, 10) }() // B -> A (sentido oposto)
	}
	wg.Wait()
	fmt.Println("total:", a.Saldo+b.Saldo) // 2000: nada se perde

	fmt.Println(Transferir(a, b, 999999)) // saldo insuficiente
}
