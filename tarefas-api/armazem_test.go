package main

import (
	"errors"
	"os"
	"path/filepath"
	"sync"
	"testing"
)

func novoCaminho(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "tarefas.json")
}

func TestArmazemFicheiroInexistenteDaListaVazia(t *testing.T) {
	a := NovoArmazem(novoCaminho(t))

	lista, err := a.Listar()
	if err != nil || len(lista) != 0 {
		t.Fatalf("lista = %v, %v; quero vazia sem erro", lista, err)
	}
}

func TestArmazemAdicionarListar(t *testing.T) {
	a := NovoArmazem(novoCaminho(t))

	t1, err := a.Adicionar("Estudar Go")
	if err != nil {
		t.Fatalf("Adicionar: %v", err)
	}
	t2, _ := a.Adicionar("Fazer API")
	if t1.ID != 1 || t2.ID != 2 {
		t.Errorf("IDs = %d, %d; quero 1, 2", t1.ID, t2.ID)
	}

	lista, _ := a.Listar()
	if len(lista) != 2 || lista[0].Titulo != "Estudar Go" || lista[1].Titulo != "Fazer API" {
		t.Errorf("lista = %+v", lista)
	}
}

func TestArmazemTituloVazio(t *testing.T) {
	caminho := novoCaminho(t)
	a := NovoArmazem(caminho)

	if _, err := a.Adicionar("  "); !errors.Is(err, ErrTituloVazio) {
		t.Errorf("erro = %v; quero ErrTituloVazio", err)
	}
	if lista, _ := a.Listar(); len(lista) != 0 {
		t.Errorf("não devia ter guardado nada: %v", lista)
	}
}

func TestArmazemIDNaoRepete(t *testing.T) {
	a := NovoArmazem(novoCaminho(t))
	a.Adicionar("a")
	a.Adicionar("b")
	if err := a.Remover(1); err != nil {
		t.Fatalf("Remover: %v", err)
	}
	c, _ := a.Adicionar("c")
	if c.ID != 3 {
		t.Errorf("ID = %d; quero 3", c.ID)
	}
}

func TestArmazemConcluirRemover(t *testing.T) {
	a := NovoArmazem(novoCaminho(t))
	a.Adicionar("a")
	a.Adicionar("b")

	if err := a.Concluir(2); err != nil {
		t.Fatalf("Concluir: %v", err)
	}
	lista, _ := a.Listar()
	if lista[0].Concluida || !lista[1].Concluida {
		t.Errorf("lista = %+v; só a 2 devia estar concluída", lista)
	}

	if err := a.Remover(1); err != nil {
		t.Fatalf("Remover: %v", err)
	}
	lista, _ = a.Listar()
	if len(lista) != 1 || lista[0].ID != 2 {
		t.Errorf("lista = %+v; quero só a tarefa 2", lista)
	}
}

func TestArmazemNaoEncontrada(t *testing.T) {
	a := NovoArmazem(novoCaminho(t))
	a.Adicionar("a")

	if err := a.Concluir(99); !errors.Is(err, ErrNaoEncontrada) {
		t.Errorf("Concluir(99) = %v; quero ErrNaoEncontrada", err)
	}
	if err := a.Remover(99); !errors.Is(err, ErrNaoEncontrada) {
		t.Errorf("Remover(99) = %v; quero ErrNaoEncontrada", err)
	}
}

// A prova de que guarda mesmo: um segundo Armazem no mesmo ficheiro
// (como um servidor reiniciado) vê as tarefas do primeiro.
func TestArmazemPersiste(t *testing.T) {
	caminho := novoCaminho(t)

	a := NovoArmazem(caminho)
	a.Adicionar("a")
	a.Adicionar("b")
	a.Concluir(1)
	a.Remover(2)

	b := NovoArmazem(caminho)
	lista, err := b.Listar()
	if err != nil {
		t.Fatalf("Listar: %v", err)
	}
	if len(lista) != 1 || lista[0].ID != 1 || !lista[0].Concluida {
		t.Errorf("lista = %+v; quero só a tarefa 1, concluída", lista)
	}
}

// Operações que falham não devem criar nem alterar o ficheiro.
func TestArmazemErroNaoCriaFicheiro(t *testing.T) {
	caminho := novoCaminho(t)
	a := NovoArmazem(caminho)

	a.Adicionar("")
	a.Concluir(1)
	a.Remover(1)

	if _, err := os.Stat(caminho); !errors.Is(err, os.ErrNotExist) {
		t.Errorf("o ficheiro não devia existir (Stat err = %v)", err)
	}
}

func TestArmazemJSONCorrompido(t *testing.T) {
	caminho := novoCaminho(t)
	if err := os.WriteFile(caminho, []byte("{isto não é json"), 0o644); err != nil {
		t.Fatal(err)
	}
	a := NovoArmazem(caminho)

	if _, err := a.Listar(); err == nil {
		t.Error("Listar devia dar erro com JSON corrompido")
	}
	if _, err := a.Adicionar("x"); err == nil {
		t.Error("Adicionar devia dar erro com JSON corrompido")
	}
	// e não pode ter destruído o ficheiro
	dados, _ := os.ReadFile(caminho)
	if string(dados) != "{isto não é json" {
		t.Errorf("o ficheiro foi alterado: %q", dados)
	}
}

// O ficheiro em disco tem de ser um array JSON válido e nunca "null".
func TestArmazemListaVaziaNoDisco(t *testing.T) {
	caminho := novoCaminho(t)
	a := NovoArmazem(caminho)
	a.Adicionar("a")
	a.Remover(1)

	dados, err := os.ReadFile(caminho)
	if err != nil {
		t.Fatal(err)
	}
	if s := string(dados); s == "null" {
		t.Errorf("ficheiro = %q; quero um array vazio []", s)
	}
}

// Com -race e IDs únicos: o Mutex tem de cobrir carregar+alterar+guardar.
func TestArmazemConcorrente(t *testing.T) {
	a := NovoArmazem(novoCaminho(t))

	const n = 50
	var wg sync.WaitGroup
	for range n {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if _, err := a.Adicionar("x"); err != nil {
				t.Errorf("Adicionar: %v", err)
			}
		}()
	}
	wg.Wait()

	lista, _ := a.Listar()
	if len(lista) != n {
		t.Fatalf("len = %d; quero %d (perdeu tarefas)", len(lista), n)
	}
	vistos := map[int]bool{}
	for _, tf := range lista {
		if vistos[tf.ID] {
			t.Errorf("ID %d repetido", tf.ID)
		}
		vistos[tf.ID] = true
	}
}

// O servidor funciona com o Armazem sem mudar uma linha.
func TestServidorComArmazem(t *testing.T) {
	var repo Repositorio = NovoArmazem(novoCaminho(t))
	repo.Adicionar("a")
	if _, err := repo.Listar(); err != nil {
		t.Fatal(err)
	}
}
