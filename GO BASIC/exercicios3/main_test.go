package main

import "testing"

func TestTrocar(t *testing.T) {
	a, b := 1, 2
	trocar(&a, &b)
	if a != 2 || b != 1 {
		t.Errorf("a=%d b=%d; esperado a=2 b=1", a, b)
	}
}

func TestDobrarTodos(t *testing.T) {
	x, y, z := 1, 5, 10
	dobrarTodos([]*int{&x, &y, &z})
	if x != 2 || y != 10 || z != 20 {
		t.Errorf("x=%d y=%d z=%d; esperado 2 10 20", x, y, z)
	}
}

func TestValorOuZero(t *testing.T) {
	n := 42
	if got := valorOuZero(&n); got != 42 {
		t.Errorf("valorOuZero(&42) = %d; esperado 42", got)
	}
	if got := valorOuZero(nil); got != 0 {
		t.Errorf("valorOuZero(nil) = %d; esperado 0", got)
	}
}

func TestContador(t *testing.T) {
	c := Contador{}
	c.Incrementar()
	c.Incrementar()
	c.Incrementar()
	if c.Valor() != 3 {
		t.Errorf("Valor() = %d; esperado 3 (o receptor de Incrementar é por ponteiro?)", c.Valor())
	}
}
