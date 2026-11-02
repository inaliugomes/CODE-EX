package main

import (
	"fmt"
	"strconv"
)

// =====================================================================
// Exercício 1 — switch com expressão
// TipoDia recebe o nome de um dia da semana e devolve:
//
//	"util"  para "segunda", "terca", "quarta", "quinta", "sexta"
//	"fim"   para "sabado", "domingo"
//	"?"     para qualquer outro valor
//
// Dica: um case aceita vários valores separados por vírgula.
// =====================================================================
func TipoDia(dia string) string {

	switch dia {
	case "segunda", "terca", "quarta", "quinta", "sexta":
		return "util"
	case "sabado", "domingo":
		return "fim"
	default:
		return "?"

	}
}

// =====================================================================
// Exercício 2 — switch sem expressão (substitui if/else if)
// Classificar devolve a classificação de uma nota de 0 a 20:
//
//	0-9   -> "negativa"
//	10-13 -> "suficiente"
//	14-17 -> "bom"
//	18-20 -> "excelente"
//	fora de 0-20 -> "invalida"
//
// Dica: `switch {` sem nada ao lado; cada case é uma condição.
// A primeira condição verdadeira ganha. Pensa na ordem dos cases.
// =====================================================================
func Classificar(nota int) string {

	switch {
	case nota >= 0 && nota <= 9:
		return "negativa"
	case nota >= 10 && nota <= 13:
		return "suficiente"
	case nota >= 14 && nota <= 17:
		return "bom"
	case nota >= 18 && nota <= 20:
		return "excelente"
	default:
		return "invalida"

	}
}

// =====================================================================
// Exercício 3 — função variádica
// Soma devolve a soma de todos os números recebidos (0 se não houver).
// Deve funcionar como Soma(), Soma(1, 2, 3) e Soma(fatia...).
// Dica: dentro da função, `nums` é um []int.
// =====================================================================
func Soma(nums ...int) int {

	var soma int

	if len(nums) == 0 {
		return 0
	}

	for _, num := range nums {

		soma += num
	}
	// TODO
	return soma
}

// =====================================================================
// Exercício 4 — variádica com parâmetro normal
// Maior devolve o maior dos números e true. Se não receber nenhum,
// devolve 0 e false.
// O parâmetro variádico tem de ser sempre o ÚLTIMO; aqui a assinatura
// já está dada: Maior(nums ...int) (int, bool).
// Atenção: não assumas que o maior é positivo (Maior(-5, -2) = -2).
// =====================================================================
func Maior(nums ...int) (int, bool) {

	if len(nums) == 0 {

		return 0, false
	}

	maior := nums[0]

	for _, num := range nums {

		if num > maior {

			maior = num

		}

	}
	// TODO
	return maior, true
}

// =====================================================================
// Exercício 5 — juntar tudo
// Descrever devolve uma frase sobre o dia e as notas recebidas:
//
//	Descrever("sabado", 12, 15) -> "fim: 2 notas"
//	Descrever("segunda")        -> "util: 0 notas"
//
// Usa TipoDia e uma variádica `notas ...int`. Sem formatação especial:
// fmt.Sprintf("%s: %d notas", ...)
// =====================================================================
func Descrever(dia string, notas ...int) string {

	tipoSeman := TipoDia(dia)
	notasStr := strconv.Itoa(len(notas))

	return fmt.Sprintf("%s: %s notas", tipoSeman, notasStr)
}

func main() {}
