package main

import (
	// Para comparar números decimais com tolerância
	"math"
	"slices"  // Tem slices.Equal, para comparar dois slices
	"testing" // O pacote de testes do Go
)

// Teste "em tabela": a forma mais comum de testar em Go.
// Defines uma lista de casos e corres o mesmo código para todos.
func TestInverter(t *testing.T) {
	casos := []struct {
		nome     string // Descrição do caso
		entrada  []int  // O que damos à função
		esperado []int  // O que ela devia devolver
	}{
		{"vários elementos", []int{1, 2, 3, 4}, []int{4, 3, 2, 1}},
		{"um elemento", []int{7}, []int{7}},
		{"slice vazio", []int{}, []int{}},
		{"números repetidos", []int{5, 5, 1}, []int{1, 5, 5}},
	}

	for _, c := range casos {
		t.Run(c.nome, func(t *testing.T) { // t.Run cria um "subteste" com nome próprio
			obtido := inverter(c.entrada)
			if !slices.Equal(obtido, c.esperado) { // Os slices não se comparam com ==
				t.Errorf("inverter(%v) = %v; esperado %v", c.entrada, obtido, c.esperado)
			}
		})
	}
}

// Testa a regra "o original não pode ser alterado"
func TestInverterNaoAlteraOriginal(t *testing.T) {
	original := []int{1, 2, 3}
	inverter(original)
	if !slices.Equal(original, []int{1, 2, 3}) {
		t.Errorf("o slice original foi alterado: %v", original)
	}
}

func TestEstatisticas(t *testing.T) {
	media, maxima, minima := estatisticas([]float64{10, 20, 30})

	// Os float64 podem ter pequenos erros de arredondamento,
	// por isso compara-se com uma margem de tolerância.
	if math.Abs(media-20) > 0.001 {
		t.Errorf("média = %v; esperado 20", media)
	}
	if maxima != 30 {
		t.Errorf("máxima = %v; esperado 30", maxima)
	}
	if minima != 10 {
		t.Errorf("mínima = %v; esperado 10", minima) // ← foi este que falhou!
	}
}

// Testa um caso-limite: o slice vazio
func TestEstatisticasSliceVazio(t *testing.T) {
	media, _, _ := estatisticas([]float64{})
	if math.IsNaN(media) {
		t.Errorf("com slice vazio, a média não devia ser NaN")
		t.Errorf("com slice vazio, a média não devia ser NaN")
	}
}
