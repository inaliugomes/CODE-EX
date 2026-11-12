package main

import (
	"errors"
	"sync"
	"testing"
)

func TestMemoriaAdicionarListar(t *testing.T) {
	m := NovaMemoria()

	lista, err := m.Listar()
	if err != nil || len(lista) != 0 {
		t.Fatalf("lista inicial = %v, %v; quero vazia sem erro", lista, err)
	}

	a, err := m.Adicionar("Estudar Go")
	if err != nil {
		t.Fatalf("Adicionar: %v", err)
	}
	b, _ := m.Adicionar("Fazer API")
	if a.ID != 1 || b.ID != 2 {
		t.Errorf("IDs = %d, %d; quero 1, 2", a.ID, b.ID)
	}
	if a.Titulo != "Estudar Go" || a.Concluida {
		t.Errorf("tarefa = %+v", a)
	}

	lista, _ = m.Listar()
	if len(lista) != 2 {
		t.Errorf("len = %d; quero 2", len(lista))
	}
}

func TestMemoriaTituloVazio(t *testing.T) {
	m := NovaMemoria()
	if _, err := m.Adicionar(""); !errors.Is(err, ErrTituloVazio) {
		t.Errorf("erro = %v; quero ErrTituloVazio", err)
	}
	if lista, _ := m.Listar(); len(lista) != 0 {
		t.Errorf("não devia ter guardado nada: %v", lista)
	}
}

// Regra herdada da CLI: ID = maior existente + 1 (nunca o tamanho da lista).
// Remover o 1 deixa [2]; com len+1 o próximo seria 2 (repetido!), com maior+1 é 3.
func TestMemoriaIDNaoRepete(t *testing.T) {
	m := NovaMemoria()
	m.Adicionar("a")
	m.Adicionar("b")
	if err := m.Remover(1); err != nil {
		t.Fatalf("Remover: %v", err)
	}
	c, _ := m.Adicionar("c")
	if c.ID != 3 {
		t.Errorf("ID = %d; quero 3", c.ID)
	}
}

func TestMemoriaConcluirRemover(t *testing.T) {
	m := NovaMemoria()
	m.Adicionar("a")

	if err := m.Concluir(1); err != nil {
		t.Fatalf("Concluir: %v", err)
	}
	lista, _ := m.Listar()
	if !lista[0].Concluida {
		t.Error("devia estar concluída")
	}

	if err := m.Concluir(99); !errors.Is(err, ErrNaoEncontrada) {
		t.Errorf("Concluir(99) = %v; quero ErrNaoEncontrada", err)
	}
	if err := m.Remover(99); !errors.Is(err, ErrNaoEncontrada) {
		t.Errorf("Remover(99) = %v; quero ErrNaoEncontrada", err)
	}
	if err := m.Remover(1); err != nil {
		t.Fatalf("Remover(1): %v", err)
	}
	if lista, _ := m.Listar(); len(lista) != 0 {
		t.Errorf("lista = %v; quero vazia", lista)
	}
}

func TestMemoriaListarDevolveCopia(t *testing.T) {
	m := NovaMemoria()
	m.Adicionar("original")
	lista, _ := m.Listar()
	lista[0].Titulo = "alterado por fora"

	outra, _ := m.Listar()
	if outra[0].Titulo != "original" {
		t.Errorf("o repositório foi alterado por fora: %q", outra[0].Titulo)
	}
}

// Corre com `go test -race`: sem Mutex, o detetor de corridas falha aqui.
func TestMemoriaConcorrente(t *testing.T) {
	m := NovaMemoria()
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			m.Adicionar("x")
			m.Listar()
		}()
	}
	wg.Wait()

	lista, _ := m.Listar()
	if len(lista) != 50 {
		t.Fatalf("len = %d; quero 50", len(lista))
	}
	vistos := map[int]bool{}
	for _, tr := range lista {
		if vistos[tr.ID] {
			t.Fatalf("ID %d repetido", tr.ID)
		}
		vistos[tr.ID] = true
	}
}
