package main

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestGuardarCarregar(t *testing.T) {
	caminho := filepath.Join(t.TempDir(), "tarefas.json")
	original := Lista{
		{ID: 1, Titulo: "Estudar Go", Concluida: true},
		{ID: 3, Titulo: "Fazer projeto"},
	}

	if err := Guardar(caminho, original); err != nil {
		t.Fatalf("Guardar: %v", err)
	}
	got, err := Carregar(caminho)
	if err != nil {
		t.Fatalf("Carregar: %v", err)
	}
	if !reflect.DeepEqual(got, original) {
		t.Errorf("got %+v, quero %+v", got, original)
	}
}

func TestCarregarInexistente(t *testing.T) {
	caminho := filepath.Join(t.TempDir(), "nao-existe.json")

	got, err := Carregar(caminho)
	if err != nil {
		t.Fatalf("erro = %v, quero nil", err)
	}
	if len(got) != 0 {
		t.Errorf("len = %d, quero 0", len(got))
	}
	if _, err := os.Stat(caminho); !os.IsNotExist(err) {
		t.Errorf("Carregar não devia criar o ficheiro")
	}
}

func TestCarregarCorrompido(t *testing.T) {
	caminho := filepath.Join(t.TempDir(), "mau.json")
	if err := os.WriteFile(caminho, []byte("{isto não é json"), 0o644); err != nil {
		t.Fatal(err)
	}

	if _, err := Carregar(caminho); err == nil {
		t.Error("JSON corrompido devia dar erro")
	}
}

func TestGuardarListaVazia(t *testing.T) {
	caminho := filepath.Join(t.TempDir(), "vazio.json")

	if err := Guardar(caminho, Lista{}); err != nil {
		t.Fatalf("Guardar: %v", err)
	}
	got, err := Carregar(caminho)
	if err != nil {
		t.Fatalf("Carregar: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("len = %d, quero 0", len(got))
	}
}

func TestGuardarDirectoriaInexistente(t *testing.T) {
	caminho := filepath.Join(t.TempDir(), "nao-existe", "tarefas.json")

	if err := Guardar(caminho, Lista{}); err == nil {
		t.Error("devia dar erro se a pasta não existe")
	}
}

// nomes devolve os nomes dos ficheiros/pastas dentro de dir.
func nomes(t *testing.T, dir string) []string {
	t.Helper()
	entradas, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	var r []string
	for _, e := range entradas {
		r = append(r, e.Name())
	}
	return r
}

// Depois de guardar com sucesso, na pasta só pode estar o ficheiro final:
// o temporário tem de desaparecer (foi renomeado).
func TestGuardarNaoDeixaTemporarios(t *testing.T) {
	dir := t.TempDir()
	caminho := filepath.Join(dir, "tarefas.json")

	for i := 0; i < 3; i++ {
		if err := Guardar(caminho, Lista{{ID: i + 1, Titulo: "x"}}); err != nil {
			t.Fatalf("Guardar: %v", err)
		}
	}

	want := []string{"tarefas.json"}
	if got := nomes(t, dir); !reflect.DeepEqual(got, want) {
		t.Errorf("ficheiros = %v, quero %v", got, want)
	}
}

// Guardar uma lista mais pequena por cima de uma grande não pode deixar
// restos do conteúdo antigo.
func TestGuardarSubstituiConteudo(t *testing.T) {
	caminho := filepath.Join(t.TempDir(), "tarefas.json")
	grande := Lista{{ID: 1, Titulo: "a"}, {ID: 2, Titulo: "b"}, {ID: 3, Titulo: "c"}}
	pequena := Lista{{ID: 1, Titulo: "a"}}

	if err := Guardar(caminho, grande); err != nil {
		t.Fatal(err)
	}
	if err := Guardar(caminho, pequena); err != nil {
		t.Fatal(err)
	}
	got, err := Carregar(caminho)
	if err != nil {
		t.Fatalf("Carregar: %v", err)
	}
	if !reflect.DeepEqual(got, pequena) {
		t.Errorf("got %+v, quero %+v", got, pequena)
	}
}

// Se o rename falhar (aqui: o destino é uma pasta não vazia), Guardar dá
// erro e tem de limpar o temporário que criou.
func TestGuardarFalhaLimpaTemporario(t *testing.T) {
	dir := t.TempDir()
	destino := filepath.Join(dir, "tarefas.json")
	if err := os.MkdirAll(filepath.Join(destino, "dentro"), 0o755); err != nil {
		t.Fatal(err)
	}

	if err := Guardar(destino, Lista{{ID: 1, Titulo: "x"}}); err == nil {
		t.Fatal("devia dar erro")
	}

	want := []string{"tarefas.json"}
	if got := nomes(t, dir); !reflect.DeepEqual(got, want) {
		t.Errorf("ficheiros = %v, quero %v (temporário não foi apagado)", got, want)
	}
}
