package main

import (
	"errors"
	"slices"
	"strings"
	"sync"
)

// Repositorio esconde ONDE as tarefas estão guardadas. O servidor HTTP só
// conhece esta interface; por isso podemos ter uma versão em memória (testes)
// e outra em ficheiro JSON (a tua lógica da CLI) sem mexer nos handlers.
type Repositorio interface {
	Listar() ([]Tarefa, error)
	Adicionar(titulo string) (Tarefa, error)
	Concluir(id int) error
	Remover(id int) error
}

// Memoria é a implementação em memória.
// TODO (Etapa 1): campos (slice de tarefas + sync.Mutex), NovaMemoria() e os 4 métodos.
//   - Listar devolve uma CÓPIA (não a slice interna): porquê?
//   - Todos os métodos protegem o acesso com o Mutex.
type Memoria struct {
	mu      sync.Mutex
	tarefas []Tarefa
}

func NovaMemoria() *Memoria { return &Memoria{} }

var errPorFazer = errors.New("por implementar")

func (m *Memoria) Listar() ([]Tarefa, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	copia := slices.Clone(m.tarefas)

	return copia, nil
}

func (m *Memoria) Adicionar(t string) (Tarefa, error) {

	if strings.TrimSpace(t) == "" { 
		return Tarefa{} , ErrTituloVazio
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	maior := 0

	for _, f := range m.tarefas {

		if f.ID > maior {
			maior = f.ID
		}

	}

	proximoID := maior + 1

	tarefa := Tarefa{ID: proximoID, Titulo: t, Concluida: false}

	m.tarefas = append(m.tarefas, tarefa)

	return tarefa, nil
}

func (m *Memoria) Concluir(id int) error { return errPorFazer }
func (m *Memoria) Remover(id int) error  { return errPorFazer }

// Garante em tempo de compilação que Memoria cumpre a interface.
var _ Repositorio = (*Memoria)(nil)
