package main

import (
	"errors"
	"fmt"
)

// =====================================================================
// Exercício 1 — ordem dos defer
// Ordem acrescenta a *log, por esta ordem de EXECUÇÃO:
//
//	"inicio", "meio", "fim-B", "fim-A"
//
// Regra: "fim-A" e "fim-B" têm de ser escritos com defer, e "fim-A"
// tem de aparecer no código ANTES de "fim-B".
// =====================================================================
func Ordem(log *[]string) {

	*log = append(*log, "inicio")
	*log = append(*log, "meio")

	defer func() {
		*log = append(*log, "fim-A")

	}()

	defer func() {

		*log = append(*log, "fim-B")
	}()

	// TODO
}

// =====================================================================
// Exercício 2 — arrumar recursos
// Recurso simula um ficheiro: Abrir/Fechar. Já está feito.
// =====================================================================
type Recurso struct {
	Aberto bool
	Fechos int
}

func (r *Recurso) Abrir()  { r.Aberto = true }
func (r *Recurso) Fechar() { r.Aberto = false; r.Fechos++ }

var ErrFalhou = errors.New("operação falhou")

// Usar abre o recurso, "trabalha" e devolve:
//   - nil se falha == false
//   - ErrFalhou se falha == true (a meio do trabalho, com return antecipado)
//
// Em AMBOS os casos o recurso tem de ficar fechado, exatamente 1 vez.
// Escreve o Fechar uma só vez, com defer, logo a seguir a Abrir.
func Usar(r *Recurso, falha bool) error {
	r.Abrir()

	defer r.Fechar()
	// TODO
	if falha == true {

		return ErrFalhou
	}
	return nil
}

// =====================================================================
// Exercício 3 — valor devolvido alterado por defer
// Dobrar devolve o dobro de x, MAS a lógica de dobrar tem de estar num
// defer: o corpo faz apenas "return x".
// Dica: precisas de um resultado com nome e de uma closure.
// =====================================================================
func Dobrar(x int) (r int) {

	defer func() {

		r = r * 2

	}()

	return x
}

// =====================================================================
// Exercício 4 — recover
// DividirSeguro devolve a/b. Se b == 0, o Go faz panic; em vez de o
// programa rebentar, devolve resultado 0 e um erro não-nil.
// Dica: recover() só funciona dentro de uma função adiada (defer), e
// para alterar o que a função devolve precisas de resultados com nome.
// =====================================================================
func DividirSeguro(a, b int) (res int, err error) {

	defer func(){
		if r := recover(); r != nil {

			res,err = 0, fmt.Errorf("panic: %v",r)
		}
	}()

	return a/b, nil
}

func main() {}
