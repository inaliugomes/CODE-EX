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

var ErrTituloVazio = errors.New("Erro Titulo Vazio")

var ErrNaoEncontrada = errors.New("Item nao encontrado")

func (t *Tarefa) Adicioanr(titulo string) (Tarefa, error) {

	return Tarefa{}, nil
}

func (t *Tarefa) Conluir(id int) error {

	return nil
}

func (t *Tarefa) Remover(id int) error {

	return nil
}

// TODO (Etapa 1): ErrTituloVazio e ErrNaoEncontrada com errors.New.

// TODO (Etapa 1): Adicionar(titulo string) (Tarefa, error)   // receptor por ponteiro
// TODO (Etapa 1): Concluir(id int) error                     // receptor por ponteiro
// TODO (Etapa 1): Remover(id int) error                      // receptor por ponteiro
