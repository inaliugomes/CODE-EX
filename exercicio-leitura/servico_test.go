package main

import (
	"errors"
	"testing"
)

func TestEncurtarGeraCodigoQuandoVazio(t *testing.T) {
	s := NovoServico(novaMemoria())
	s.gerar = func() string { return "abc123" }

	l, err := s.Encurtar("https://go.dev", "")
	if err != nil {
		t.Fatal(err)
	}
	if l.Codigo != "abc123" {
		t.Errorf("codigo = %q, quero abc123", l.Codigo)
	}
}

func TestEncurtarRejeitaURLInvalido(t *testing.T) {
	s := NovoServico(novaMemoria())
	_, err := s.Encurtar("isto não é um url", "")
	if !errors.Is(err, ErrURLInvalido) {
		t.Errorf("err = %v, quero ErrURLInvalido", err)
	}
}

func TestEncurtarCodigoRepetido(t *testing.T) {
	s := NovoServico(novaMemoria())
	s.Encurtar("https://go.dev", "go")
	_, err := s.Encurtar("https://gobyexample.com", "go")
	if !errors.Is(err, ErrCodigoExiste) {
		t.Errorf("err = %v, quero ErrCodigoExiste", err)
	}
}

func TestSeguirContaVisitas(t *testing.T) {
	s := NovoServico(novaMemoria())
	s.Encurtar("https://go.dev", "go")
	s.Seguir("go")
	s.Seguir("go")

	l, _ := s.Estatisticas("go")
	if l.Visitas != 2 {
		t.Errorf("visitas = %d, quero 2", l.Visitas)
	}
}

func TestSeguirInexistente(t *testing.T) {
	s := NovoServico(novaMemoria())
	_, err := s.Seguir("nada")
	if !errors.Is(err, ErrNaoExiste) {
		t.Errorf("err = %v, quero ErrNaoExiste", err)
	}
}
