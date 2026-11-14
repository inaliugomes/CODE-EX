package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

type Servidor struct {
	repo Repositorio

	mux *http.ServeMux
}

func NovoServidor(repo Repositorio) *Servidor {

	s := &Servidor{repo: repo, mux: http.NewServeMux()}

	s.mux.HandleFunc("GET /tarefas", s.listar)
	s.mux.HandleFunc("POST /tarefas", s.adicionar)
	s.mux.HandleFunc("PATCH /tarefas/{id}/concluir", s.concluir)

	return s
}

func (s *Servidor) Handler() http.Handler { return s.mux }

func (s *Servidor) listar(w http.ResponseWriter, r *http.Request) {

	tarefas, err := s.repo.Listar()

	if err != nil {

		http.Error(w, "Erro Interno do servidor , tente mas tarde ", http.StatusInternalServerError)
		return

	}

	if tarefas == nil {
		tarefas = []Tarefa{}
	}

	w.Header().Set("Content-type", "application/json")
	errencode := json.NewEncoder(w).Encode(tarefas)

	if errencode != nil {

		http.Error(w, "Erro Interno do servidor , tente mas tarde ", http.StatusInternalServerError)
		return
	}
}

func (s *Servidor) adicionar(w http.ResponseWriter, r *http.Request) {

	var entrada struct {
		Titulo string `json:"titulo"`
	}

	err := json.NewDecoder(r.Body).Decode(&entrada)

	if err != nil {

		http.Error(w, "JSON inválido ou titulo vazio", http.StatusBadRequest)
		return

	}

	tarefa, errTarefa := s.repo.Adicionar(entrada.Titulo)

	if errTarefa != nil {

		if errors.Is(errTarefa, ErrTituloVazio) {
			http.Error(w, "JSON inválido ou titulo vazio", http.StatusBadRequest)
			return
		} else {

			http.Error(w, "Erro do servidor , por favor tente mas tarde", http.StatusInternalServerError)
			return

		}
	}

	w.Header().Set("Content-type", "application/json")

	w.WriteHeader(201)
	json.NewEncoder(w).Encode(tarefa)

}

func (s *Servidor) concluir(w http.ResponseWriter, r *http.Request) {

	id := r.PathValue("id")

	idValue, errValue := strconv.Atoi(id)

	if errValue != nil {

		http.Error(w, "Erro do servidor , valida o seu id", http.StatusBadRequest)
		return
	}

	errConcluir := s.repo.Concluir(idValue)

	if errConcluir != nil {
		if errors.Is(errConcluir, ErrNaoEncontrada) {

			http.Error(w, "tarefa com esse id nao encontrado", http.StatusNotFound)
			return
		} else {

			http.Error(w, "Error interno do servidor", http.StatusInternalServerError)
			return

		}
	}

	w.WriteHeader(http.StatusNoContent)

}
