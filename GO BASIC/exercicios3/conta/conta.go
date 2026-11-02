// Package conta implementa uma conta bancária simples.
package conta

// =====================================================================
// Conta tem campos privados (minúscula): só este pacote lhes acede.
// =====================================================================
type Conta struct {
	titular string
	saldo   float64
}

// =====================================================================
// Nova cria uma conta com saldo 0 e devolve um PONTEIRO para ela.
// =====================================================================
func Nova(titular string) *Conta {

	// TODO
	return &Conta{titular: titular, saldo: 0}
}

// =====================================================================
// Saldo devolve o saldo atual.
// =====================================================================
func (c *Conta) Saldo() float64 {
	// TODO
	return c.saldo
}

// =====================================================================
// Depositar soma valor ao saldo.
// Se valor <= 0, não altera nada e devolve ErrValorInvalido.
// =====================================================================
func (c *Conta) Depositar(valor float64) error {
	// TODO

	if valor > 0 {
		c.saldo += valor
		return nil
	}
	return ErrValorInvalido
}

// =====================================================================
// Levantar subtrai valor ao saldo.
//   - valor <= 0          -> ErrValorInvalido
//   - valor > saldo       -> ErrSaldoInsuficiente
//
// Em caso de erro, o saldo não muda.
// =====================================================================
func (c *Conta) Levantar(valor float64) error {

	if valor <= 0 {
		return ErrValorInvalido
	}

	if valor > c.saldo {

		return ErrSaldoInsuficiente
	}

	c.saldo -= valor
	// TODO
	return nil
}
