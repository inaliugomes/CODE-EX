package main

import "testing"

func TestTipoDia(t *testing.T) {
	casos := []struct{ dia, esperado string }{
		{"segunda", "util"}, {"sexta", "util"}, {"quarta", "util"},
		{"sabado", "fim"}, {"domingo", "fim"},
		{"", "?"}, {"Segunda", "?"}, {"feriado", "?"},
	}
	for _, c := range casos {
		if got := TipoDia(c.dia); got != c.esperado {
			t.Errorf("TipoDia(%q) = %q; esperado %q", c.dia, got, c.esperado)
		}
	}
}

func TestClassificar(t *testing.T) {
	casos := []struct {
		nota     int
		esperado string
	}{
		{-1, "invalida"}, {0, "negativa"}, {9, "negativa"},
		{10, "suficiente"}, {13, "suficiente"},
		{14, "bom"}, {17, "bom"},
		{18, "excelente"}, {20, "excelente"}, {21, "invalida"},
	}
	for _, c := range casos {
		if got := Classificar(c.nota); got != c.esperado {
			t.Errorf("Classificar(%d) = %q; esperado %q", c.nota, got, c.esperado)
		}
	}
}

func TestSoma(t *testing.T) {
	if got := Soma(); got != 0 {
		t.Errorf("Soma() = %d; esperado 0", got)
	}
	if got := Soma(1, 2, 3); got != 6 {
		t.Errorf("Soma(1,2,3) = %d; esperado 6", got)
	}
	fatia := []int{10, 20, -5}
	if got := Soma(fatia...); got != 25 {
		t.Errorf("Soma(fatia...) = %d; esperado 25", got)
	}
}

func TestMaior(t *testing.T) {
	if got, ok := Maior(); got != 0 || ok {
		t.Errorf("Maior() = %d, %v; esperado 0, false", got, ok)
	}
	if got, ok := Maior(3, 9, 4); got != 9 || !ok {
		t.Errorf("Maior(3,9,4) = %d, %v; esperado 9, true", got, ok)
	}
	if got, ok := Maior(-5, -2, -9); got != -2 || !ok {
		t.Errorf("Maior(-5,-2,-9) = %d, %v; esperado -2, true", got, ok)
	}
	if got, ok := Maior(0); got != 0 || !ok {
		t.Errorf("Maior(0) = %d, %v; esperado 0, true", got, ok)
	}
}

func TestDescrever(t *testing.T) {
	if got := Descrever("sabado", 12, 15); got != "fim: 2 notas" {
		t.Errorf("Descrever(sabado,12,15) = %q", got)
	}
	if got := Descrever("segunda"); got != "util: 0 notas" {
		t.Errorf("Descrever(segunda) = %q", got)
	}
	if got := Descrever("xpto", 1); got != "?: 1 notas" {
		t.Errorf("Descrever(xpto,1) = %q", got)
	}
}
