package main

import (
	"errors"
	"fmt"
)

// =====================================================================
// EXERCÍCIO 1 — Contar palavras
// ---------------------------------------------------------------------
// Recebes um slice de palavras. Devolve um map que diz quantas vezes
// cada palavra aparece.
//
// Exemplo: ["go", "é", "go"]  ->  map[go:2 é:1]
//
// Dica: num map, ler uma chave que não existe devolve o valor zero
// (0 para int), por isso contagem[p]++ funciona sempre.
// =====================================================================
func contarPalavras(palavras []string) map[string]int {

	contador := make(map[string]int)

	for _, p := range palavras {

		contador[p]++

	}
	// TODO
	return contador
}

// =====================================================================
// EXERCÍCIO 2 — Inverter um map
// ---------------------------------------------------------------------
// Recebe um map[string]int e devolve um map[int]string com as chaves
// e os valores trocados.
//
// Exemplo: {"a": 1, "b": 2}  ->  {1: "a", 2: "b"}
//
// Podes assumir que os valores originais são todos diferentes.
// =====================================================================
func inverterMap(m map[string]int) map[int]string {

	chavetrocado := make(map[int]string)

	for key, value := range m {

		chavetrocado[value] = key
	}
	// TODO
	return chavetrocado
}

// =====================================================================
// EXERCÍCIO 3 — Erro de divisão
// ---------------------------------------------------------------------
// Escreve "dividir(a, b float64) (float64, error)".
// Se b for 0, devolve o erro ErrDivisaoPorZero (já declarado abaixo).
// Caso contrário devolve o resultado e nil.
// =====================================================================
var ErrDivisaoPorZero = errors.New("divisão por zero")

func dividir(a, b float64) (float64, error) {

	if b == 0 {

		return 0, ErrDivisaoPorZero
	}
	result := a / b

	return result, nil
}

// =====================================================================
// EXERCÍCIO 4 — Procurar utilizador (erro com contexto)
// ---------------------------------------------------------------------
// "utilizadores" é um map de id -> nome.
// Escreve "procurar(id int) (string, error)":
//   - se o id existir, devolve o nome e nil
//   - se não existir, devolve um erro que EMBRULHA ErrNaoEncontrado
//     e inclui o id na mensagem, usando fmt.Errorf com %w
//
// Exemplo de mensagem: "utilizador 99: não encontrado"
//
// Dica: o teste usa errors.Is(err, ErrNaoEncontrado). Isso só funciona
// se usares %w (e não %v) no fmt.Errorf.
// Dica 2: lê um map com dois valores: nome, ok := utilizadores[id]
// =====================================================================
var ErrNaoEncontrado = errors.New("não encontrado")

var utilizadores = map[int]string{
	1: "Ana",
	2: "Bruno",
	3: "Carla",
}

// Procurar dentro de um dicionario usando for talvez nao seja a melhor maneira , o que devemos fazer é passar um valor chave , caso nao existe
// Validar o resultado atraves de um processo de erro.
func procurar(id int) (string, error) {

	nome, ok := utilizadores[id]

	if !ok {
		return "", fmt.Errorf("utilizador %d: %w", id, ErrNaoEncontrado)
	}

	return nome, nil
}

// =====================================================================
// EXERCÍCIO 5 — Tipo de erro próprio (desafio)
// ---------------------------------------------------------------------
// Cria um tipo "ErroIdade" (struct com o campo Idade int) que
// implementa a interface error (método Error() string).
//
// Escreve "validarIdade(idade int) error":
//   - idade < 0 ou > 150  ->  devolve ErroIdade{Idade: idade}
//   - caso contrário      ->  devolve nil
//
// A mensagem pode ser: "idade inválida: <valor>"
//
// O teste usa errors.As para extrair o ErroIdade do erro.
// =====================================================================

type ErroIdade struct {
	Idade int
}

func (e ErroIdade) Error() string {

	idadestr := fmt.Sprintf("%d", e.Idade)

	idadeError := "idade inválida:" + idadestr

	return idadeError
}

func validarIdade(idade int) error {
	// TODO

	if idade < 0 || idade > 150 {

		return ErroIdade{Idade: idade}
	}
	return nil
}

func main() {
	fmt.Println(contarPalavras([]string{"go", "é", "go"}))
}
