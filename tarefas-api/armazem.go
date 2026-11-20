package main

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// Armazem é o Repositorio que guarda as tarefas num ficheiro JSON.
// Cada método faz Lock -> carregar -> alterar -> guardar: o Mutex cobre a
// sequência toda, senão dois pedidos em paralelo perdiam dados.
type Armazem struct {
	mu      sync.Mutex
	caminho string
}

var _ Repositorio = (*Armazem)(nil)

func NovoArmazem(caminho string) *Armazem { return &Armazem{caminho: caminho} }

func (a *Armazem) Listar() ([]Tarefa, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	return carregar(a.caminho)
}

func (a *Armazem) Adicionar(titulo string) (Tarefa, error) {
	if strings.TrimSpace(titulo) == "" {
		return Tarefa{}, ErrTituloVazio
	}

	a.mu.Lock()
	defer a.mu.Unlock()

	lista, err := carregar(a.caminho)
	if err != nil {
		return Tarefa{}, err
	}

	// ID = maior existente + 1 (nunca len, senão repete depois de remover).
	maior := 0
	for _, t := range lista {
		if t.ID > maior {
			maior = t.ID
		}
	}

	nova := Tarefa{ID: maior + 1, Titulo: titulo}
	if err := guardar(a.caminho, append(lista, nova)); err != nil {
		return Tarefa{}, err
	}

	return nova, nil
}

func (a *Armazem) Concluir(id int) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	lista, err := carregar(a.caminho)
	if err != nil {
		return err
	}

	for i := range lista {
		if lista[i].ID == id {
			lista[i].Concluida = true // por índice: o range daria uma cópia
			return guardar(a.caminho, lista)
		}
	}

	return ErrNaoEncontrada // não encontrada: o ficheiro fica intacto
}

func (a *Armazem) Remover(id int) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	lista, err := carregar(a.caminho)
	if err != nil {
		return err
	}

	for i := range lista {
		if lista[i].ID == id {
			return guardar(a.caminho, append(lista[:i], lista[i+1:]...))
		}
	}

	return ErrNaoEncontrada
}

// carregar lê a lista do ficheiro. Ficheiro inexistente = lista vazia (1.º arranque).
func carregar(caminho string) ([]Tarefa, error) {
	dados, err := os.ReadFile(caminho)
	if errors.Is(err, os.ErrNotExist) {
		return []Tarefa{}, nil
	}
	if err != nil {
		return nil, err
	}

	var lista []Tarefa
	if err := json.Unmarshal(dados, &lista); err != nil {
		return nil, err
	}

	return lista, nil
}

// guardar escreve de forma segura: ficheiro temporário na mesma pasta + Rename
// (atómico). Se algo falhar a meio, o ficheiro original não fica corrompido.
func guardar(caminho string, lista []Tarefa) error {
	if lista == nil {
		lista = []Tarefa{} // nil viraria "null" no JSON
	}

	dados, err := json.MarshalIndent(lista, "", " ")
	if err != nil {
		return err
	}

	tmp, err := os.CreateTemp(filepath.Dir(caminho), "tarefas-*.tmp")
	if err != nil {
		return err
	}

	if _, err := tmp.Write(dados); err != nil {
		tmp.Close()
		os.Remove(tmp.Name())
		return err
	}

	// O erro do Close conta: verificar antes do Rename.
	if err := tmp.Close(); err != nil {
		os.Remove(tmp.Name())
		return err
	}

	if err := os.Rename(tmp.Name(), caminho); err != nil {
		os.Remove(tmp.Name())
		return err
	}

	return nil
}
