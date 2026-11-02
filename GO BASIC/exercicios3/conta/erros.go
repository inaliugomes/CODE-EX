package conta

import "errors"

// Erros exportados: quem usa o pacote pode compará-los com errors.Is.
var (
	ErrValorInvalido     = errors.New("valor inválido")
	ErrSaldoInsuficiente = errors.New("saldo insuficiente")
)
