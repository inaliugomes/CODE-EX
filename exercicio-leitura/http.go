package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

func NovoHandler(s *Servico) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /links", criar(s))
	mux.HandleFunc("GET /links/{codigo}", estatisticas(s))
	mux.HandleFunc("GET /{codigo}", seguir(s))
	return registar(mux)
}

func criar(s *Servico) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var entrada struct {
			Destino string `json:"destino"`
			Codigo  string `json:"codigo"`
		}
		if err := json.NewDecoder(r.Body).Decode(&entrada); err != nil {
			responderErro(w, http.StatusBadRequest, "json inválido")
			return
		}
		l, err := s.Encurtar(entrada.Destino, entrada.Codigo)
		if err != nil {
			responderErro(w, estado(err), err.Error())
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(l)
	}
}

func seguir(s *Servico) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		destino, err := s.Seguir(r.PathValue("codigo"))
		if err != nil {
			responderErro(w, estado(err), err.Error())
			return
		}
		http.Redirect(w, r, destino, http.StatusFound)
	}
}

func estatisticas(s *Servico) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l, err := s.Estatisticas(r.PathValue("codigo"))
		if err != nil {
			responderErro(w, estado(err), err.Error())
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(l)
	}
}

func estado(err error) int {
	switch {
	case errors.Is(err, ErrURLInvalido):
		return http.StatusBadRequest
	case errors.Is(err, ErrCodigoExiste):
		return http.StatusConflict
	case errors.Is(err, ErrNaoExiste):
		return http.StatusNotFound
	}
	return http.StatusInternalServerError
}

func responderErro(w http.ResponseWriter, codigo int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(codigo)
	json.NewEncoder(w).Encode(map[string]string{"erro": msg})
}

func registar(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
