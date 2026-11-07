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
