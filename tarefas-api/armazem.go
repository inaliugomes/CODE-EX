package main

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sync"
)

// Armazem é o Repositorio que guarda as tarefas num ficheiro JSON.
//
// TODO (Etapa 5):
//  1. Copia Carregar e Guardar de ../tarefas-cli/armazem.go para este ficheiro,
//     trocando o tipo Lista por []Tarefa (e escolhe nomes: ficam como funções
//     soltas ou métodos privados).
//  2. Cada método faz: Lock -> defer Unlock -> carregar -> alterar -> guardar.
//     O Mutex cobre a sequência TODA (senão dois pedidos perdem dados).
//  3. Só guardar quando a operação teve sucesso (erro => ficheiro intacto).
//  4. Lista vazia grava-se como [] (nunca null): atenção à slice nil.
type Armazem struct {
	mu      sync.Mutex
	caminho string
}

var errPorFazerArmazem = errors.New("por fazer")

func NovoArmazem(caminho string) *Armazem { return &Armazem{caminho: caminho} }

func (a *Armazem) Listar() ([]Tarefa, error) {

	a.mu.Lock()
	defer a.mu.Unlock()

	listaObjeto, err := Carregar(a.caminho)

	if err != nil {

		return nil, err

	}

	return listaObjeto, nil

}

func (a *Armazem) Adicionar(t string) (Tarefa, error) { return Tarefa{}, errPorFazerArmazem }
func (a *Armazem) Concluir(id int) error              { return errPorFazerArmazem }
func (a *Armazem) Remover(id int) error               { return errPorFazerArmazem }

var _ Repositorio = (*Armazem)(nil)

func Carregar(caminho string) ([]Tarefa, error) {

	f, err := os.ReadFile(caminho)

	if errors.Is(err, os.ErrNotExist) {

		return []Tarefa{}, nil
	}

	if err != nil {

		return []Tarefa{}, err
	}

	var list []Tarefa

	jsonObjeto := json.Unmarshal(f, &list)

	if jsonObjeto != nil {

		return []Tarefa{}, jsonObjeto
	}

	return list, nil

}

func Guardar(caminho string, l []Tarefa) error {

	jsonObject, err := json.MarshalIndent(l, "", " ")

	if err != nil {

		return err

	}

	caminhoDir := filepath.Dir(caminho)

	tmp, errPath := os.CreateTemp(caminhoDir, "tarefas-*.tmp")

	if errPath != nil {

		return errPath

	}

	_, errWriteJSon := tmp.Write(jsonObject)

	if errWriteJSon != nil {

		tmp.Close()

		os.Remove(tmp.Name())

		return errWriteJSon

	}

	errClose := tmp.Close()

	if errClose != nil {

		os.Remove(tmp.Name())

		return errClose

	}

	errRename := os.Rename(tmp.Name(), caminho)

	if errRename != nil {

		os.Remove(tmp.Name())

		return errRename

	}

	return nil

}
