package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func pedir(t *testing.T, repo Repositorio, metodo, caminho string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(metodo, caminho, nil)
	rec := httptest.NewRecorder()
	NovoServidor(repo).Handler().ServeHTTP(rec, req)
	return rec
}

func TestListarVazioDevolveArrayVazio(t *testing.T) {
	rec := pedir(t, NovaMemoria(), http.MethodGet, "/tarefas")

	if rec.Code != http.StatusOK {
		t.Errorf("código = %d, esperado 200", rec.Code)
	}
	if corpo := strings.TrimSpace(rec.Body.String()); corpo != "[]" {
		t.Errorf("corpo = %q, esperado %q (nunca null)", corpo, "[]")
	}
}

func TestListarDevolveTarefasEmJSON(t *testing.T) {
	repo := NovaMemoria()
	repo.Adicionar("Estudar Go")
	repo.Adicionar("Escrever testes")

	rec := pedir(t, repo, http.MethodGet, "/tarefas")

	if rec.Code != http.StatusOK {
		t.Fatalf("código = %d, esperado 200", rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); !strings.HasPrefix(ct, "application/json") {
		t.Errorf("Content-Type = %q, esperado application/json", ct)
	}
	var got []Tarefa
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("corpo não é JSON válido: %v", err)
	}
	if len(got) != 2 || got[0].Titulo != "Estudar Go" || got[1].Titulo != "Escrever testes" {
		t.Errorf("tarefas = %+v", got)
	}
}

func TestListarUsaCamposJSONEmMinusculas(t *testing.T) {
	repo := NovaMemoria()
	repo.Adicionar("Estudar Go")

	rec := pedir(t, repo, http.MethodGet, "/tarefas")

	esperado := `[{"id":1,"titulo":"Estudar Go","concluida":false}]`
	if corpo := strings.TrimSpace(rec.Body.String()); corpo != esperado {
		t.Errorf("corpo = %s, esperado %s", corpo, esperado)
	}
}

type repoComErro struct{ Repositorio }

func (repoComErro) Listar() ([]Tarefa, error) { return nil, errors.New("falhou") }

func TestListarErroDoRepositorioDa500(t *testing.T) {
	rec := pedir(t, repoComErro{NovaMemoria()}, http.MethodGet, "/tarefas")

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("código = %d, esperado 500", rec.Code)
	}
}

func TestRotaDesconhecidaDa404(t *testing.T) {
	rec := pedir(t, NovaMemoria(), http.MethodGet, "/nada")

	if rec.Code != http.StatusNotFound {
		t.Errorf("código = %d, esperado 404", rec.Code)
	}
}
