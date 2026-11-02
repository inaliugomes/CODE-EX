package conta

import (
	"errors"
	"testing"
)

func TestNovaContaTemSaldoZero(t *testing.T) {
	c := Nova("Ana")
	if c == nil {
		t.Fatal("Nova devolveu nil")
	}
	if c.Saldo() != 0 {
		t.Errorf("Saldo() = %v; esperado 0", c.Saldo())
	}
}

func TestDepositar(t *testing.T) {
	c := Nova("Ana")
	if err := c.Depositar(100); err != nil {
		t.Fatalf("Depositar(100) devolveu %v", err)
	}
	if c.Saldo() != 100 {
		t.Errorf("Saldo() = %v; esperado 100", c.Saldo())
	}

	for _, v := range []float64{0, -10} {
		err := c.Depositar(v)
		if !errors.Is(err, ErrValorInvalido) {
			t.Errorf("Depositar(%v) = %v; esperado ErrValorInvalido", v, err)
		}
	}
	if c.Saldo() != 100 {
		t.Errorf("depósitos inválidos alteraram o saldo: %v", c.Saldo())
	}
}

func TestLevantar(t *testing.T) {
	c := Nova("Ana")
	_ = c.Depositar(100)

	if err := c.Levantar(30); err != nil {
		t.Fatalf("Levantar(30) devolveu %v", err)
	}
	if c.Saldo() != 70 {
		t.Errorf("Saldo() = %v; esperado 70", c.Saldo())
	}

	if err := c.Levantar(0); !errors.Is(err, ErrValorInvalido) {
		t.Errorf("Levantar(0) = %v; esperado ErrValorInvalido", err)
	}

	if err := c.Levantar(1000); !errors.Is(err, ErrSaldoInsuficiente) {
		t.Errorf("Levantar(1000) = %v; esperado ErrSaldoInsuficiente", err)
	}
	if err := c.Levantar(-5); !errors.Is(err, ErrValorInvalido) {
		t.Errorf("Levantar(-5) = %v; esperado ErrValorInvalido", err)
	}
	if c.Saldo() != 70 {
		t.Errorf("operações inválidas alteraram o saldo: %v", c.Saldo())
	}
}
