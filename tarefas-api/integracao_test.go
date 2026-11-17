package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
)

// Teste de ponta a ponta: servidor real (porta aleatória), cliente HTTP real,
// repositório em ficheiro. Ao contrário do Recorder, aqui passa pela rede.
func TestFluxoCompletoPorHTTP(t *testing.T) {
	caminho := filepath.Join(t.TempDir(), "tarefas.json")
	srv := httptest.NewServer(NovoServidor(NovoArmazem(caminho)).Handler())
	defer srv.Close()

	// 1. POST cria
	resp, err := http.Post(srv.URL+"/tarefas", "application/json",
		strings.NewReader(`{"titulo":"Estudar Go"}`))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("POST: código %d, esperado 201", resp.StatusCode)
	}
	var criada Tarefa
	if err := json.NewDecoder(resp.Body).Decode(&criada); err != nil {
		t.Fatal(err)
	}
	if criada.ID != 1 || criada.Titulo != "Estudar Go" {
		t.Errorf("tarefa criada = %+v", criada)
	}

	// 2. GET lista
	listar := func() []Tarefa {
		r, err := http.Get(srv.URL + "/tarefas")
		if err != nil {
			t.Fatal(err)
		}
		defer r.Body.Close()
		var l []Tarefa
		if err := json.NewDecoder(r.Body).Decode(&l); err != nil {
			t.Fatal(err)
		}
		return l
	}
	if l := listar(); len(l) != 1 || l[0].Titulo != "Estudar Go" {
		t.Fatalf("GET: lista = %+v", l)
	}

	// 3. PATCH conclui, 4. DELETE remove (o http.Client não tem atalho: NewRequest)
	for _, c := range []struct{ metodo, caminho string }{
		{http.MethodPatch, "/tarefas/1/concluir"},
		{http.MethodDelete, "/tarefas/1"},
	} {
		req, _ := http.NewRequest(c.metodo, srv.URL+c.caminho, nil)
		r, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		r.Body.Close()
		if r.StatusCode != http.StatusNoContent {
			t.Errorf("%s %s: código %d, esperado 204", c.metodo, c.caminho, r.StatusCode)
		}
	}

	if l := listar(); len(l) != 0 {
		t.Errorf("depois do DELETE: lista = %+v, esperada vazia", l)
	}
}

// Os dados têm de sobreviver a um "reinício": novo servidor, mesmo ficheiro.
func TestPersistenciaEntreServidores(t *testing.T) {
	caminho := filepath.Join(t.TempDir(), "tarefas.json")

	srv1 := httptest.NewServer(NovoServidor(NovoArmazem(caminho)).Handler())
	resp, err := http.Post(srv1.URL+"/tarefas", "application/json",
		strings.NewReader(`{"titulo":"Persistir"}`))
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	srv1.Close()

	srv2 := httptest.NewServer(NovoServidor(NovoArmazem(caminho)).Handler())
	defer srv2.Close()
	r, err := http.Get(srv2.URL + "/tarefas")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Body.Close()
	var l []Tarefa
	if err := json.NewDecoder(r.Body).Decode(&l); err != nil {
		t.Fatal(err)
	}
	if len(l) != 1 || l[0].Titulo != "Persistir" {
		t.Errorf("depois do reinício: lista = %+v", l)
	}
}
