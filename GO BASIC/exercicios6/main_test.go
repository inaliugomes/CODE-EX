package main

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestGuardarCarregar(t *testing.T) {
	caminho := filepath.Join(t.TempDir(), "tarefas.json")
	original := []Tarefa{
		{ID: 1, Texto: "estudar Go", Feita: true},
		{ID: 2, Texto: "fazer exercícios", Feita: false},
	}
	if err := Guardar(caminho, original); err != nil {
		t.Fatalf("Guardar: %v", err)
	}

	dados, err := os.ReadFile(caminho)
	if err != nil {
		t.Fatalf("ficheiro não foi criado: %v", err)
	}
	if !json.Valid(dados) {
		t.Errorf("o ficheiro não tem JSON válido: %q", dados)
	}
	if len(dados) > 0 && !reflect.DeepEqual(dados[:2], []byte("[\n")) {
		t.Errorf("esperava JSON indentado (começa com \"[\\n\"), veio %q", dados)
	}

	lidas, err := Carregar(caminho)
	if err != nil {
		t.Fatalf("Carregar: %v", err)
	}
	if !reflect.DeepEqual(lidas, original) {
		t.Errorf("Carregar = %+v; esperado %+v", lidas, original)
	}
}

func TestCarregarInexistente(t *testing.T) {
	caminho := filepath.Join(t.TempDir(), "nao-existe.json")
	lidas, err := Carregar(caminho)
	if err != nil {
		t.Errorf("ficheiro inexistente não é erro; veio %v", err)
	}
	if len(lidas) != 0 {
		t.Errorf("esperava lista vazia; veio %+v", lidas)
	}
}

func TestCarregarJSONInvalido(t *testing.T) {
	caminho := filepath.Join(t.TempDir(), "mau.json")
	if err := os.WriteFile(caminho, []byte("{isto não é json"), 0644); err != nil {
		t.Fatal(err)
	}
	if _, err := Carregar(caminho); err == nil {
		t.Error("JSON inválido devia dar erro")
	}
}

func TestGravar(t *testing.T) {
	caminho := filepath.Join(t.TempDir(), "a.txt")
	if err := Gravar(caminho, "primeiro"); err != nil {
		t.Fatalf("Gravar: %v", err)
	}
	if err := Gravar(caminho, "segundo"); err != nil {
		t.Fatalf("Gravar: %v", err)
	}
	dados, _ := os.ReadFile(caminho)
	if string(dados) != "segundo" {
		t.Errorf("Gravar devia substituir; ficheiro tem %q", dados)
	}
}

func TestGravarErro(t *testing.T) {
	// pasta que não existe: o Create falha
	caminho := filepath.Join(t.TempDir(), "nao-existe", "a.txt")
	if err := Gravar(caminho, "x"); err == nil {
		t.Error("Gravar numa pasta inexistente devia dar erro")
	}
}

func TestAcrescentar(t *testing.T) {
	caminho := filepath.Join(t.TempDir(), "log.txt")
	for _, l := range []string{"um", "dois", "três"} {
		if err := Acrescentar(caminho, l); err != nil {
			t.Fatalf("Acrescentar(%q): %v", l, err)
		}
	}
	dados, _ := os.ReadFile(caminho)
	if string(dados) != "um\ndois\ntrês\n" {
		t.Errorf("ficheiro = %q; esperado %q", dados, "um\ndois\ntrês\n")
	}
}

func TestAcrescentarNaoApaga(t *testing.T) {
	caminho := filepath.Join(t.TempDir(), "log.txt")
	if err := os.WriteFile(caminho, []byte("antigo\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := Acrescentar(caminho, "novo"); err != nil {
		t.Fatal(err)
	}
	dados, _ := os.ReadFile(caminho)
	if string(dados) != "antigo\nnovo\n" {
		t.Errorf("ficheiro = %q; esperado %q", dados, "antigo\nnovo\n")
	}
}

func TestLerLinhas(t *testing.T) {
	caminho := filepath.Join(t.TempDir(), "l.txt")
	if err := os.WriteFile(caminho, []byte("a\nb\nc\n"), 0644); err != nil {
		t.Fatal(err)
	}
	linhas, err := LerLinhas(caminho)
	if err != nil {
		t.Fatalf("LerLinhas: %v", err)
	}
	if !reflect.DeepEqual(linhas, []string{"a", "b", "c"}) {
		t.Errorf("LerLinhas = %q; esperado [a b c]", linhas)
	}
}

func TestLerLinhasVazioEInexistente(t *testing.T) {
	dir := t.TempDir()
	vazio := filepath.Join(dir, "vazio.txt")
	if err := os.WriteFile(vazio, nil, 0644); err != nil {
		t.Fatal(err)
	}
	linhas, err := LerLinhas(vazio)
	if err != nil || len(linhas) != 0 {
		t.Errorf("ficheiro vazio: linhas=%q err=%v; esperado vazio e nil", linhas, err)
	}

	_, err = LerLinhas(filepath.Join(dir, "nao-existe.txt"))
	if !errors.Is(err, os.ErrNotExist) {
		t.Errorf("inexistente: err=%v; esperado errors.Is(err, os.ErrNotExist)", err)
	}
}
