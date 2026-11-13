package main

import (
	"encoding/json"
	"net/http"
)

// Servidor liga o Repositorio ao HTTP.
// TODO (Etapa 2 em diante): campo com o Repositorio e Handler() http.Handler
// com um http.ServeMux e as rotas do PLANO.md.
type Servidor struct {
	repo Repositorio

	mux *http.ServeMux
}

func NovoServidor(repo Repositorio) *Servidor {

	s := &Servidor{repo: repo, mux: http.NewServeMux()}

	s.mux.HandleFunc("GET /tarefas", s.listar)

	return s
}

func (s *Servidor) Handler() http.Handler { return s.mux }

func (s *Servidor) listar(w http.ResponseWriter, r *http.Request) {

	tarefas, err := s.repo.Listar()

	if err != nil {

		http.Error(w, "....", http.StatusInternalServerError)
		return

	}

	if tarefas == nil {
		tarefas = []Tarefa{}
	}

	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(tarefas)
}
