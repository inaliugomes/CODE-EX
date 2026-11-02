package main

import "fmt"

// =====================================================================
// EXERCÍCIO 1 — Temperaturas da semana
// ---------------------------------------------------------------------
// Recebes um slice com as temperaturas de 7 dias (float64).
// Usa um ciclo "for ... range" para calcular:
//   - a MÉDIA das temperaturas
//   - a temperatura MÁXIMA
//   - a temperatura MÍNIMA
//
// A função deve devolver os três valores, por esta ordem.
//
// Dica: para a máxima e a mínima, começa com o primeiro elemento do
// slice (temperaturas[0]) e vai comparando com os restantes.
// =====================================================================
func estatisticas(temperaturas []float64) (media, maxima, minima float64) {

	if len(temperaturas) == 0 {

		return 0, 0, 0
	}

	maxima = temperaturas[0]
	minima = temperaturas[0]

	soma := 0.0

	for _, t := range temperaturas {

		soma += t
		if t > maxima {
			maxima = t
		}

		if t < minima {
			minima = t
		}
	}

	media = soma / float64(len(temperaturas))

	return media, maxima, minima
}

// =====================================================================
// EXERCÍCIO 2 — Inverter um slice
// ---------------------------------------------------------------------
// Escreve a função "inverter" que recebe um slice de inteiros e
// devolve um NOVO slice com os elementos por ordem inversa.
//
// Exemplo: [1 2 3 4]  ->  [4 3 2 1]
//
// Regra: o slice ORIGINAL não pode ser alterado!
//
// Dica: lembra-te de que os slices partilham memória. Pensa em como
// usar make() para criar um slice novo com o tamanho certo.
// =====================================================================
func inverter(s []int) []int {
	// ✏️ ESCREVE AQUI A TUA LÓGICA

	var sliceInvertido []int

	for i := len(s) - 1; i >= 0; i-- {

		sliceInvertido = append(sliceInvertido, s[i])
	}

	return sliceInvertido // substitui "nil" pelo teu slice invertido
}

// =====================================================================
// NÃO PRECISAS DE ALTERAR O CÓDIGO ABAIXO.
// Serve apenas para testar as tuas funções.
// =====================================================================
func main() {
	fmt.Println("===== Exercício 1 =====")
	semana := []float64{18.5, 21.0, 19.2, 24.8, 22.1, 17.3, 20.0}
	media, maxima, minima := estatisticas(semana)
	fmt.Printf("Média: %.2f | Máxima: %.1f | Mínima: %.1f\n", media, maxima, minima)
	fmt.Println("Esperado: Média: 20.41 | Máxima: 24.8 | Mínima: 17.3")

	fmt.Println()
	fmt.Println("===== Exercício 2 =====")
	original := []int{1, 2, 3, 4}
	invertido := inverter(original)
	fmt.Println("Original: ", original, " (deve continuar [1 2 3 4])")
	fmt.Println("Invertido:", invertido, " (deve ser [4 3 2 1])")
}
