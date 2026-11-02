package main

import (
	"exercicios3/conta"
	"fmt"
)

// =====================================================================
// PARTE 1 — PONTEIROS
// Lembrete:  &x  = endereço de x      *p  = valor para onde p aponta
// =====================================================================

// =====================================================================
// EXERCÍCIO 1 — Trocar dois valores
// ---------------------------------------------------------------------
// Recebe dois ponteiros para int e troca os valores a que apontam.
//
// Exemplo: a=1, b=2  ->  depois de trocar(&a, &b): a=2, b=1
// =====================================================================
func trocar(a, b *int) {

	temp := *b

	*b = *a
	*a = temp

}

// =====================================================================
// EXERCÍCIO 2 — Dobrar todos
// ---------------------------------------------------------------------
// Recebe um slice de ponteiros para int e dobra o valor de cada um.
//
// Dica: usa *p = *p * 2 dentro de um for ... range.
// =====================================================================
func dobrarTodos(numeros []*int) {
	// TODO

	for i := 0; i < len(numeros); i++ {
		//Aqui o que estou a dizer é o seguiente , vai neste indereço memoria
		//Pega o valor que esta nele e atualizar o valor para * 2
		*numeros[i] = *numeros[i] * 2
	}
}

// =====================================================================
// EXERCÍCIO 3 — Ponteiro nil
// ---------------------------------------------------------------------
// Devolve o valor para onde p aponta. Se p for nil, devolve 0.
//
// Dica: desreferenciar um ponteiro nil provoca panic!
// =====================================================================
func valorOuZero(p *int) int {
	// TODO

	if p == nil {

		return 0
	} else {

		return *p
	}

}

// =====================================================================
// EXERCÍCIO 4 — Receptor por valor vs por ponteiro
// ---------------------------------------------------------------------
// Completa os dois métodos de Contador:
//   - Incrementar: soma 1 ao campo N (tem de alterar o original!)
//   - Valor: devolve N (só lê, por isso pode usar receptor por valor)
//
// Pergunta para pensares: o que acontece se Incrementar tiver
// receptor por valor (c Contador) em vez de (c *Contador)?
// =====================================================================
type Contador struct {
	N int
}

func (c *Contador) Incrementar() {
	// TODO: corrige o receptor e implementa
	c.N++
}

func (c Contador) Valor() int {

	return c.N

}

// =====================================================================
// PARTE 2 — PACOTES (ver a pasta conta/)
// ---------------------------------------------------------------------
// O pacote "conta" está dividido em dois ficheiros:
//   conta/conta.go   -> o tipo Conta e os seus métodos
//   conta/erros.go   -> os erros do pacote
// Completa os TODO e corre:  go test ./...
//
// Repara: nomes com letra MAIÚSCULA são exportados (visíveis fora do
// pacote); com minúscula são privados. Por isso os campos de Conta são
// privados e só se mexe neles pelos métodos.
// =====================================================================

func main() {
	x, y := 1, 2
	trocar(&x, &y)
	fmt.Println(x, y)

	c := conta.Nova("Ana")
	_ = c.Depositar(100)
	fmt.Println(c.Saldo())
}
