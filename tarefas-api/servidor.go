package main

import "net/http"

// Servidor liga o Repositorio ao HTTP.
// TODO (Etapa 2 em diante): campo com o Repositorio e Handler() http.Handler
// com um http.ServeMux e as rotas do PLANO.md.
type Servidor struct {
}

func NovoServidor(repo Repositorio) *Servidor { return &Servidor{} }

func (s *Servidor) Handler() http.Handler { return http.NewServeMux() }
