package main

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"
)

// corre executa os comandos em sequência, sobre o mesmo ficheiro,
// e devolve a saída e o erro do ÚLTIMO comando.
func corre(t *testing.T, caminho string, comandos ...[]string) (string, error) {
	t.Helper()
	var saida bytes.Buffer
	var err error
	for _, args := range comandos {
		saida.Reset()
		err = executar(args, caminho, &saida)
	}
	return saida.String(), err
}

func TestExecutarFluxo(t *testing.T) {
	caminho := filepath.Join(t.TempDir(), "tarefas.json")

	out, err := corre(t, caminho, []string{"add", "Estudar Go"})
	if err != nil || out != "Tarefa 1 criada.\n" {
		t.Errorf("add: out=%q err=%v", out, err)
	}

	out, err = corre(t, caminho, []string{"list"})
	if err != nil || out != "[ ] 1  Estudar Go\n" {
		t.Errorf("list: out=%q err=%v", out, err)
	}

	out, err = corre(t, caminho, []string{"done", "1"})
	if err != nil || out != "Tarefa 1 concluída.\n" {
		t.Errorf("done: out=%q err=%v", out, err)
	}

	out, err = corre(t, caminho, []string{"list"})
	if err != nil || out != "[x] 1  Estudar Go\n" {
		t.Errorf("list após done: out=%q err=%v", out, err)
	}

	out, err = corre(t, caminho, []string{"remove", "1"})
	if err != nil || out != "Tarefa 1 removida.\n" {
		t.Errorf("remove: out=%q err=%v", out, err)
	}

	out, err = corre(t, caminho, []string{"list"})
	if err != nil || out != "" {
		t.Errorf("list vazia: out=%q err=%v", out, err)
	}
}

func TestExecutarPersiste(t *testing.T) {
	caminho := filepath.Join(t.TempDir(), "tarefas.json")

	corre(t, caminho, []string{"add", "a"}, []string{"add", "b"})
	out, _ := corre(t, caminho, []string{"list"})

	want := "[ ] 1  a\n[ ] 2  b\n"
	if out != want {
		t.Errorf("out = %q, quero %q (cada executar deve carregar e guardar)", out, want)
	}
}

func TestExecutarErros(t *testing.T) {
	tabela := []struct {
		nome string
		args []string
	}{
		{"sem argumentos", nil},
		{"comando desconhecido", []string{"voar"}},
		{"add sem título", []string{"add"}},
		{"add com título vazio", []string{"add", ""}},
		{"done sem ID", []string{"done"}},
		{"done inexistente", []string{"done", "99"}},
		{"remove sem ID", []string{"remove"}},
		{"remove inexistente", []string{"remove", "99"}},
	}
	for _, tc := range tabela {
		t.Run(tc.nome, func(t *testing.T) {
			caminho := filepath.Join(t.TempDir(), "tarefas.json")
			if _, err := corre(t, caminho, tc.args); err == nil {
				t.Error("devia dar erro")
			}
		})
	}
}

func TestExecutarAjuda(t *testing.T) {
	var saida bytes.Buffer
	caminho := filepath.Join(t.TempDir(), "tarefas.json")

	_ = executar(nil, caminho, &saida)
	if saida.Len() == 0 {
		t.Error("sem argumentos devia escrever a ajuda")
	}
}

// ID que não é um número deve dar o erro "id inválido" (e não "não
// encontrada"), sem panic e sem alterar a lista.
func TestExecutarIDInvalido(t *testing.T) {
	ids := []string{"abc", "1.5", "1a", " ", "99999999999999999999"}
	for _, cmd := range []string{"done", "remove"} {
		for _, id := range ids {
			t.Run(cmd+" "+id, func(t *testing.T) {
				caminho := filepath.Join(t.TempDir(), "tarefas.json")
				corre(t, caminho, []string{"add", "a"})

				_, err := corre(t, caminho, []string{cmd, id})
				if err == nil || !strings.Contains(err.Error(), "id inválido") {
					t.Errorf("err = %v, quero um erro com %q", err, "id inválido")
				}

				out, _ := corre(t, caminho, []string{"list"})
				if out != "[ ] 1  a\n" {
					t.Errorf("a lista mudou: %q", out)
				}
			})
		}
	}
}
