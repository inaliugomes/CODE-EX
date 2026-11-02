package main

import (
	"errors"
	"maps"
	"testing"
)

func TestContarPalavras(t *testing.T) {
	obtido := contarPalavras([]string{"go", "é", "go", "go", "fácil"})
	esperado := map[string]int{"go": 3, "é": 1, "fácil": 1}
	if !maps.Equal(obtido, esperado) {
		t.Errorf("obtido %v; esperado %v", obtido, esperado)
	}
}

func TestContarPalavrasVazio(t *testing.T) {
	obtido := contarPalavras(nil)
	if len(obtido) != 0 {
		t.Errorf("obtido %v; esperado map vazio", obtido)
	}
}

func TestInverterMap(t *testing.T) {
	obtido := inverterMap(map[string]int{"a": 1, "b": 2, "c": 3})
	esperado := map[int]string{1: "a", 2: "b", 3: "c"}
	if !maps.Equal(obtido, esperado) {
		t.Errorf("obtido %v; esperado %v", obtido, esperado)
	}
}

func TestDividir(t *testing.T) {
	r, err := dividir(10, 4)
	if err != nil || r != 2.5 {
		t.Errorf("dividir(10, 4) = %v, %v; esperado 2.5, nil", r, err)
	}

	_, err = dividir(1, 0)
	if !errors.Is(err, ErrDivisaoPorZero) {
		t.Errorf("dividir(1, 0) devolveu %v; esperado ErrDivisaoPorZero", err)
	}
}

func TestProcurar(t *testing.T) {
	nome, err := procurar(2)
	if err != nil || nome != "Bruno" {
		t.Errorf("procurar(2) = %q, %v; esperado \"Bruno\", nil", nome, err)
	}

	_, err = procurar(99)
	if !errors.Is(err, ErrNaoEncontrado) {
		t.Fatalf("procurar(99) devolveu %v; esperado erro que embrulhe ErrNaoEncontrado (usa %%w)", err)
	}
	if err == ErrNaoEncontrado {
		t.Errorf("o erro devia ter contexto (o id), não ser o ErrNaoEncontrado simples")
	}
}

func TestValidarIdade(t *testing.T) {
	if err := validarIdade(30); err != nil {
		t.Errorf("validarIdade(30) = %v; esperado nil", err)
	}

	err := validarIdade(-5)
	var e ErroIdade
	if !errors.As(err, &e) {
		t.Fatalf("validarIdade(-5) devolveu %v; esperado ErroIdade", err)
	}
	if e.Idade != -5 {
		t.Errorf("e.Idade = %d; esperado -5", e.Idade)
	}

	if err := validarIdade(200); err == nil {
		t.Errorf("validarIdade(200) devia devolver erro")
	}
}
