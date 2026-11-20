package main

import (
	"math/rand/v2"
	"net/url"
)

const letras = "abcdefghijklmnopqrstuvwxyz0123456789"

type Servico struct {
	armazem Armazem
	gerar   func() string
}

func NovoServico(a Armazem) *Servico {
	return &Servico{armazem: a, gerar: codigoAleatorio}
}

func codigoAleatorio() string {
	b := make([]byte, 6)
	for i := range b {
		b[i] = letras[rand.IntN(len(letras))]
	}
	return string(b)
}

// Encurtar cria um link. Se codigo vier vazio, gera um.
func (s *Servico) Encurtar(destino, codigo string) (Link, error) {
	u, err := url.ParseRequestURI(destino)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") {
		return Link{}, ErrURLInvalido
	}
	if codigo == "" {
		codigo = s.gerar()
	}
	l := Link{Codigo: codigo, Destino: destino}
	if err := s.armazem.Criar(l); err != nil {
		return Link{}, err
	}
	return l, nil
}

// Seguir conta uma visita e devolve o destino.
func (s *Servico) Seguir(codigo string) (string, error) {
	l, err := s.armazem.RegistarVisita(codigo)
	if err != nil {
		return "", err
	}
	return l.Destino, nil
}

func (s *Servico) Estatisticas(codigo string) (Link, error) {
	return s.armazem.Procurar(codigo)
}
