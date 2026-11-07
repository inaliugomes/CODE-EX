package main

import (
	"errors"
	"testing"
)

func TestAdicionar(t *testing.T) {
	tabela := []struct {
		nome    string
		inicial Lista
		titulo  string
		wantID  int
		wantErr error
		wantLen int
	}{
		{"lista vazia dá ID 1", Lista{}, "Estudar Go", 1, nil, 1},
		{"maior ID + 1", Lista{{ID: 1, Titulo: "a"}, {ID: 2, Titulo: "b"}}, "c", 3, nil, 3},
		{"ID não reutiliza buraco", Lista{{ID: 1, Titulo: "a"}, {ID: 3, Titulo: "c"}}, "d", 4, nil, 3},
		{"título vazio", Lista{}, "", 0, ErrTituloVazio, 0},
	}
	for _, tc := range tabela {
		t.Run(tc.nome, func(t *testing.T) {
			l := tc.inicial
			got, err := l.Adicionar(tc.titulo)
			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("erro = %v, quero %v", err, tc.wantErr)
			}
			if err != nil {
				if len(l) != tc.wantLen {
					t.Errorf("len = %d, quero %d (lista não devia mudar)", len(l), tc.wantLen)
				}
				return
			}
			if got.ID != tc.wantID || got.Titulo != tc.titulo || got.Concluida {
				t.Errorf("tarefa = %+v", got)
			}
			if len(l) != tc.wantLen {
				t.Errorf("len = %d, quero %d", len(l), tc.wantLen)
			}
		})
	}
}

func TestConcluir(t *testing.T) {
	tabela := []struct {
		nome    string
		inicial Lista
		id      int
		wantErr error
	}{
		{"existente", Lista{{ID: 1, Titulo: "a"}}, 1, nil},
		{"já concluída não é erro", Lista{{ID: 1, Titulo: "a", Concluida: true}}, 1, nil},
		{"inexistente", Lista{{ID: 1, Titulo: "a"}}, 99, ErrNaoEncontrada},
		{"lista vazia", Lista{}, 1, ErrNaoEncontrada},
	}
	for _, tc := range tabela {
		t.Run(tc.nome, func(t *testing.T) {
			l := tc.inicial
			err := l.Concluir(tc.id)
			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("erro = %v, quero %v", err, tc.wantErr)
			}
			if err == nil && !l[0].Concluida {
				t.Errorf("tarefa %d devia estar concluída", tc.id)
			}
		})
	}
}

func TestRemover(t *testing.T) {
	tabela := []struct {
		nome    string
		inicial Lista
		id      int
		wantIDs []int
		wantErr error
	}{
		{"do meio", Lista{{ID: 1}, {ID: 2}, {ID: 3}}, 2, []int{1, 3}, nil},
		{"o primeiro", Lista{{ID: 1}, {ID: 2}}, 1, []int{2}, nil},
		{"o último", Lista{{ID: 1}, {ID: 2}}, 2, []int{1}, nil},
		{"único elemento", Lista{{ID: 1}}, 1, []int{}, nil},
		{"inexistente", Lista{{ID: 1}}, 99, []int{1}, ErrNaoEncontrada},
		{"lista vazia", Lista{}, 1, []int{}, ErrNaoEncontrada},
	}
	for _, tc := range tabela {
		t.Run(tc.nome, func(t *testing.T) {
			l := tc.inicial
			err := l.Remover(tc.id)
			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("erro = %v, quero %v", err, tc.wantErr)
			}
			if len(l) != len(tc.wantIDs) {
				t.Fatalf("len = %d, quero %d", len(l), len(tc.wantIDs))
			}
			for i, id := range tc.wantIDs {
				if l[i].ID != id {
					t.Errorf("l[%d].ID = %d, quero %d (ordem deve manter-se)", i, l[i].ID, id)
				}
			}
		})
	}
}
