package main

import (
	"reflect"
	"sort"
	"testing"
	"time"
)

// comTimeout falha o teste em vez de ficar preso para sempre se a função
// bloquear (deadlock, channel nunca fechado, etc.).
func comTimeout(t *testing.T, f func()) {
	t.Helper()
	feito := make(chan struct{})
	go func() {
		f()
		close(feito)
	}()
	select {
	case <-feito:
	case <-time.After(2 * time.Second):
		t.Fatal("timeout: a função ficou bloqueada (deadlock ou channel por fechar?)")
	}
}

func TestSomarParalelo(t *testing.T) {
	casos := []struct {
		nome string
		in   []int
		out  int
	}{
		{"nil", nil, 0},
		{"um elemento", []int{7}, 7},
		{"tamanho ímpar", []int{1, 2, 3}, 6},
		{"dez", []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, 55},
		{"negativos", []int{-1, -2, 5}, 2},
	}
	for _, c := range casos {
		t.Run(c.nome, func(t *testing.T) {
			comTimeout(t, func() {
				if got := SomarParalelo(c.in); got != c.out {
					t.Errorf("SomarParalelo(%v) = %d; esperado %d", c.in, got, c.out)
				}
			})
		})
	}
}

func TestGerar(t *testing.T) {
	comTimeout(t, func() {
		got := []int{}
		for v := range Gerar(5) {
			got = append(got, v)
		}
		if !reflect.DeepEqual(got, []int{1, 2, 3, 4, 5}) {
			t.Errorf("Gerar(5) = %v; esperado [1 2 3 4 5]", got)
		}
	})
}

func TestGerarZero(t *testing.T) {
	comTimeout(t, func() {
		n := 0
		for range Gerar(0) {
			n++
		}
		if n != 0 {
			t.Errorf("Gerar(0) emitiu %d valores; esperado 0", n)
		}
	})
}

func TestQuadradosParalelo(t *testing.T) {
	comTimeout(t, func() {
		got := QuadradosParalelo([]int{1, 2, 3, 4, -5})
		if !reflect.DeepEqual(got, []int{1, 4, 9, 16, 25}) {
			t.Errorf("QuadradosParalelo = %v; esperado [1 4 9 16 25]", got)
		}
		if got := QuadradosParalelo(nil); len(got) != 0 {
			t.Errorf("QuadradosParalelo(nil) = %v; esperado vazio", got)
		}
	})
}

func TestDobrar(t *testing.T) {
	comTimeout(t, func() {
		in := make(chan int, 3)
		in <- 1
		in <- 2
		in <- 3
		close(in)

		got := []int{}
		for v := range Dobrar(in) {
			got = append(got, v)
		}
		if !reflect.DeepEqual(got, []int{2, 4, 6}) {
			t.Errorf("Dobrar = %v; esperado [2 4 6]", got)
		}
	})
}

func TestPipelineGerarDobrar(t *testing.T) {
	comTimeout(t, func() {
		got := []int{}
		for v := range Dobrar(Gerar(4)) {
			got = append(got, v)
		}
		if !reflect.DeepEqual(got, []int{2, 4, 6, 8}) {
			t.Errorf("Dobrar(Gerar(4)) = %v; esperado [2 4 6 8]", got)
		}
	})
}

func TestTrabalhadores(t *testing.T) {
	for _, n := range []int{1, 3, 20} {
		comTimeout(t, func() {
			tarefas := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
			got := Trabalhadores(tarefas, n)
			sort.Ints(got)
			esperado := []int{10, 20, 30, 40, 50, 60, 70, 80, 90, 100}
			if !reflect.DeepEqual(got, esperado) {
				t.Errorf("Trabalhadores(n=%d) = %v; esperado %v", n, got, esperado)
			}
		})
	}
}

func TestTrabalhadoresSemTarefas(t *testing.T) {
	comTimeout(t, func() {
		if got := Trabalhadores(nil, 3); len(got) != 0 {
			t.Errorf("Trabalhadores(nil) = %v; esperado vazio", got)
		}
	})
}
