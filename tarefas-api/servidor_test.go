package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
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

func pedirComCorpo(t *testing.T, repo Repositorio, metodo, caminho, corpo string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(metodo, caminho, strings.NewReader(corpo))
	rec := httptest.NewRecorder()
	NovoServidor(repo).Handler().ServeHTTP(rec, req)
	return rec
}

func TestAdicionarDevolve201ETarefa(t *testing.T) {
	repo := NovaMemoria()
	rec := pedirComCorpo(t, repo, http.MethodPost, "/tarefas", `{"titulo":"Estudar Go"}`)

	if rec.Code != http.StatusCreated {
		t.Fatalf("código = %d, esperado 201", rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); !strings.HasPrefix(ct, "application/json") {
		t.Errorf("Content-Type = %q, esperado application/json", ct)
	}
	esperado := `{"id":1,"titulo":"Estudar Go","concluida":false}`
	if corpo := strings.TrimSpace(rec.Body.String()); corpo != esperado {
		t.Errorf("corpo = %s, esperado %s", corpo, esperado)
	}
	if lista, _ := repo.Listar(); len(lista) != 1 {
		t.Errorf("o repositório devia ter 1 tarefa, tem %d", len(lista))
	}
}

func TestAdicionarJSONInvalidoDa400(t *testing.T) {
	repo := NovaMemoria()
	for _, corpo := range []string{`isto não é json`, ``, `{"titulo":`} {
		rec := pedirComCorpo(t, repo, http.MethodPost, "/tarefas", corpo)
		if rec.Code != http.StatusBadRequest {
			t.Errorf("corpo %q: código = %d, esperado 400", corpo, rec.Code)
		}
	}
	if lista, _ := repo.Listar(); len(lista) != 0 {
		t.Errorf("não devia ter criado nada, tem %d", len(lista))
	}
}

func TestAdicionarTituloVazioDa400(t *testing.T) {
	repo := NovaMemoria()
	for _, corpo := range []string{`{"titulo":""}`, `{}`, `{"titulo":"   "}`} {
		rec := pedirComCorpo(t, repo, http.MethodPost, "/tarefas", corpo)
		if rec.Code != http.StatusBadRequest {
			t.Errorf("corpo %q: código = %d, esperado 400", corpo, rec.Code)
		}
	}
}

type repoAdicionarErro struct{ Repositorio }

func (repoAdicionarErro) Adicionar(string) (Tarefa, error) {
	return Tarefa{}, errors.New("falhou")
}

func TestAdicionarErroInesperadoDa500SemCorpoDeSucesso(t *testing.T) {
	rec := pedirComCorpo(t, repoAdicionarErro{NovaMemoria()}, http.MethodPost, "/tarefas", `{"titulo":"x"}`)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("código = %d, esperado 500", rec.Code)
	}
	if strings.Contains(rec.Body.String(), `"id"`) {
		t.Errorf("o corpo não devia ter uma tarefa: %s", rec.Body.String())
	}
}

func TestAdicionarEmParalelo(t *testing.T) {
	repo := NovaMemoria()
	srv := NovoServidor(repo).Handler()

	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			req := httptest.NewRequest(http.MethodPost, "/tarefas", strings.NewReader(`{"titulo":"x"}`))
			srv.ServeHTTP(httptest.NewRecorder(), req)
		}()
	}
	wg.Wait()

	if lista, _ := repo.Listar(); len(lista) != 50 {
		t.Errorf("tarefas = %d, esperado 50", len(lista))
	}
}

// --- Etapa 4: PATCH /tarefas/{id}/concluir e DELETE /tarefas/{id} ---

func TestConcluirDevolve204EMarcaConcluida(t *testing.T) {
	repo := NovaMemoria()
	repo.Adicionar("Estudar Go")

	rec := pedir(t, repo, http.MethodPatch, "/tarefas/1/concluir")

	if rec.Code != http.StatusNoContent {
		t.Fatalf("código = %d, esperado 204", rec.Code)
	}
	if rec.Body.Len() != 0 {
		t.Errorf("204 não deve ter corpo, tem %q", rec.Body.String())
	}
	if lista, _ := repo.Listar(); len(lista) != 1 || !lista[0].Concluida {
		t.Errorf("a tarefa devia estar concluída: %+v", lista)
	}
}

func TestConcluirInexistenteDa404(t *testing.T) {
	rec := pedir(t, NovaMemoria(), http.MethodPatch, "/tarefas/99/concluir")

	if rec.Code != http.StatusNotFound {
		t.Errorf("código = %d, esperado 404", rec.Code)
	}
}

func TestConcluirIDInvalidoDa400(t *testing.T) {
	for _, id := range []string{"abc", "1.5", "-"} {
		rec := pedir(t, NovaMemoria(), http.MethodPatch, "/tarefas/"+id+"/concluir")
		if rec.Code != http.StatusBadRequest {
			t.Errorf("id %q: código = %d, esperado 400", id, rec.Code)
		}
	}
}

func TestRemoverDevolve204ERemove(t *testing.T) {
	repo := NovaMemoria()
	repo.Adicionar("A")
	repo.Adicionar("B")

	rec := pedir(t, repo, http.MethodDelete, "/tarefas/1")

	if rec.Code != http.StatusNoContent {
		t.Fatalf("código = %d, esperado 204", rec.Code)
	}
	if rec.Body.Len() != 0 {
		t.Errorf("204 não deve ter corpo, tem %q", rec.Body.String())
	}
	if lista, _ := repo.Listar(); len(lista) != 1 || lista[0].Titulo != "B" {
		t.Errorf("devia restar só a tarefa B: %+v", lista)
	}
}

func TestRemoverInexistenteDa404(t *testing.T) {
	rec := pedir(t, NovaMemoria(), http.MethodDelete, "/tarefas/99")

	if rec.Code != http.StatusNotFound {
		t.Errorf("código = %d, esperado 404", rec.Code)
	}
}

func TestRemoverIDInvalidoDa400(t *testing.T) {
	rec := pedir(t, NovaMemoria(), http.MethodDelete, "/tarefas/abc")

	if rec.Code != http.StatusBadRequest {
		t.Errorf("código = %d, esperado 400", rec.Code)
	}
}

type repoConcluirRemoverErro struct{ Repositorio }

func (repoConcluirRemoverErro) Concluir(int) error { return errors.New("falhou") }
func (repoConcluirRemoverErro) Remover(int) error  { return errors.New("falhou") }

func TestConcluirEErroInesperadoDa500(t *testing.T) {
	repo := repoConcluirRemoverErro{NovaMemoria()}

	if rec := pedir(t, repo, http.MethodPatch, "/tarefas/1/concluir"); rec.Code != http.StatusInternalServerError {
		t.Errorf("concluir: código = %d, esperado 500", rec.Code)
	}
	if rec := pedir(t, repo, http.MethodDelete, "/tarefas/1"); rec.Code != http.StatusInternalServerError {
		t.Errorf("remover: código = %d, esperado 500", rec.Code)
	}
}

func TestMetodoErradoNaRotaDa405(t *testing.T) {
	rec := pedir(t, NovaMemoria(), http.MethodPut, "/tarefas/1/concluir")

	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("código = %d, esperado 405", rec.Code)
	}
}
