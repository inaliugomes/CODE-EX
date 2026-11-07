package main

import (
	"errors"
)

// DECISÕES DA ETAPA 0 (responde tu, por escrito, aqui):
//
// 1. Como se gera o ID? E porque não usar o tamanho da lista?
//    R:ID DEVEM SER UNICOS , eu sei sei GO TEM ALGUMA FUNCAO de gerar id , mas acredito uma das for mas u
//
// 2. Concluir duas vezes a mesma tarefa: é erro ou não faz nada?
//    R:
//
// 3. O que acontece se o ficheiro JSON não existir?
//    R:

type Tarefa struct {
	ID        int    `json:"id"`
	Titulo    string `json:"titulo"`
	Concluida bool   `json:"concluida"`
}

type Lista []Tarefa

var ErrTituloVazio = errors.New("erro titulo vazio")

var ErrNaoEncontrada = errors.New("tem nao encontrado")

func (l *Lista) Adicionar(titulo string) (Tarefa, error) {

	if titulo == "" {

		return Tarefa{}, ErrTituloVazio
	}

	maior := 0

	for _, valor := range *l {

		if valor.ID > maior {

			maior = valor.ID
		}

	}

	nextID := maior + 1

	tarefa := Tarefa{ID: nextID, Titulo: titulo, Concluida: false}

	*l = append(*l, tarefa)

	return tarefa, nil

}

func (l *Lista) Concluir(id int) error {

	if id < 0 {

		return ErrNaoEncontrada
	}

	for i, tarefa := range *l {

		if tarefa.ID == id {

			(*l)[i].Concluida = true

			return nil

		}

	}

	return ErrNaoEncontrada

}

func (l *Lista) Remover(id int) error {

	if id < 0 {

		return ErrNaoEncontrada
	}

	for i, tarefa := range *l {

		if tarefa.ID == id {

			*l = append((*l)[:i], (*l)[i+1:]...)
			return nil
		}
	}

	return ErrNaoEncontrada

}
