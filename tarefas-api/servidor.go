package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
)

// Servidor liga as rotas HTTP ao Repositorio. Só conhece a interface.
type Servidor struct {
	repo Repositorio
	mux  *http.ServeMux
}

func NovoServidor(repo Repositorio) *Servidor {
	s := &Servidor{repo: repo, mux: http.NewServeMux()}

	s.mux.HandleFunc("GET /tarefas", s.listar)
	s.mux.HandleFunc("POST /tarefas", s.adicionar)
	s.mux.HandleFunc("PATCH /tarefas/{id}/concluir", s.concluir)
	s.mux.HandleFunc("DELETE /tarefas/{id}", s.remover)

	return s
}

func (s *Servidor) Handler() http.Handler { return s.mux }

// escreverJSON responde com o código e o corpo em JSON.
// A ordem importa: cabeçalho -> WriteHeader -> corpo.
func escreverJSON(w http.ResponseWriter, codigo int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(codigo)
	// Depois do WriteHeader já não dá para mudar o código: só registamos o erro.
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Println("erro a escrever resposta:", err)
	}
}

// escreverErro responde com {"erro":"mensagem"}.
func escreverErro(w http.ResponseWriter, codigo int, mensagem string) {
	escreverJSON(w, codigo, map[string]string{"erro": mensagem})
}

// lerID lê {id} do caminho. Erro se não for um número.
func lerID(r *http.Request) (int, error) {
	return strconv.Atoi(r.PathValue("id"))
}

func (s *Servidor) listar(w http.ResponseWriter, r *http.Request) {
	tarefas, err := s.repo.Listar()
	if err != nil {
		escreverErro(w, http.StatusInternalServerError, "erro interno do servidor")
		return
	}

	// nil viraria "null" em JSON; o contrato é sempre uma lista.
	if tarefas == nil {
		tarefas = []Tarefa{}
	}

	escreverJSON(w, http.StatusOK, tarefas)
}

func (s *Servidor) adicionar(w http.ResponseWriter, r *http.Request) {
	// Local ao handler: uma variável global seria uma data race entre pedidos.
	var entrada struct {
		Titulo string `json:"titulo"`
	}

	if err := json.NewDecoder(r.Body).Decode(&entrada); err != nil {
		escreverErro(w, http.StatusBadRequest, "JSON inválido")
		return
	}

	tarefa, err := s.repo.Adicionar(entrada.Titulo)
	if errors.Is(err, ErrTituloVazio) {
		escreverErro(w, http.StatusBadRequest, "título vazio")
		return
	}
	if err != nil {
		escreverErro(w, http.StatusInternalServerError, "erro interno do servidor")
		return
	}

	escreverJSON(w, http.StatusCreated, tarefa)
}

func (s *Servidor) concluir(w http.ResponseWriter, r *http.Request) {
	id, err := lerID(r)
	if err != nil {
		escreverErro(w, http.StatusBadRequest, "id inválido")
		return
	}

	s.responderSemCorpo(w, s.repo.Concluir(id))
}

func (s *Servidor) remover(w http.ResponseWriter, r *http.Request) {
	id, err := lerID(r)
	if err != nil {
		escreverErro(w, http.StatusBadRequest, "id inválido")
		return
	}

	s.responderSemCorpo(w, s.repo.Remover(id))
}

// responderSemCorpo traduz o resultado de Concluir/Remover: 204, 404 ou 500.
// O repositório não sabe de HTTP; é aqui que o erro vira código.
func (s *Servidor) responderSemCorpo(w http.ResponseWriter, err error) {
	if errors.Is(err, ErrNaoEncontrada) {
		escreverErro(w, http.StatusNotFound, "tarefa não encontrada")
		return
	}
	if err != nil {
		escreverErro(w, http.StatusInternalServerError, "erro interno do servidor")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
