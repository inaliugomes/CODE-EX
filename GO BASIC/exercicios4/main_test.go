package main

import (
	"errors"
	"reflect"
	"testing"
)

func TestOrdem(t *testing.T) {
	var log []string
	Ordem(&log)
	esperado := []string{"inicio", "meio", "fim-B", "fim-A"}
	if !reflect.DeepEqual(log, esperado) {
		t.Errorf("Ordem = %v; esperado %v", log, esperado)
	}
}

func TestUsar(t *testing.T) {
	for _, falha := range []bool{false, true} {
		r := &Recurso{}
		err := Usar(r, falha)
		if falha && !errors.Is(err, ErrFalhou) {
			t.Errorf("Usar(falha=true) = %v; esperado ErrFalhou", err)
		}
		if !falha && err != nil {
			t.Errorf("Usar(falha=false) = %v; esperado nil", err)
		}
		if r.Aberto {
			t.Errorf("falha=%v: recurso ficou aberto", falha)
		}
		if r.Fechos != 1 {
			t.Errorf("falha=%v: Fechos = %d; esperado 1", falha, r.Fechos)
		}
	}
}

func TestDobrar(t *testing.T) {
	for _, c := range []struct{ in, out int }{{0, 0}, {3, 6}, {-4, -8}} {
		if got := Dobrar(c.in); got != c.out {
			t.Errorf("Dobrar(%d) = %d; esperado %d", c.in, got, c.out)
		}
	}
}

func TestDividirSeguro(t *testing.T) {
	if res, err := DividirSeguro(10, 2); res != 5 || err != nil {
		t.Errorf("DividirSeguro(10, 2) = %d, %v; esperado 5, nil", res, err)
	}
	if res, err := DividirSeguro(7, 0); res != 0 || err == nil {
		t.Errorf("DividirSeguro(7, 0) = %d, %v; esperado 0 e erro", res, err)
	}
}
