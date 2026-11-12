package main

import "errors"

// Reaproveitado da CLI. Podes copiar o teu código de ../tarefas-cli/tarefa.go
// ou reescrever à tua maneira.

type Tarefa struct {
	ID        int    `json:"id"`
	Titulo    string `json:"titulo"`
	Concluida bool   `json:"concluida"`
}

var (
	ErrTituloVazio   = errors.New("título vazio")
	ErrNaoEncontrada = errors.New("tarefa não encontrada")
)
