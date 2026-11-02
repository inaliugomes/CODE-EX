package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

// =====================================================================
// Tarefa já está feita. Repara nas maiúsculas e nas tags json.
// =====================================================================
type Tarefa struct {
	ID    int    `json:"id"`
	Texto string `json:"texto"`
	Feita bool   `json:"feita"`
}

// =====================================================================
// Exercício 1 — Guardar (JSON + os.WriteFile)
// Guardar converte as tarefas para JSON com indentação de 2 espaços
// (json.MarshalIndent) e escreve-as em caminho com permissões 0644.
// Se algo falhar, devolve o erro embrulhado com contexto (%w).
// =====================================================================
func Guardar(caminho string, tarefas []Tarefa) error {

	tarefasJSon, err := json.MarshalIndent(tarefas, "", "  ")

	if err != nil {

		return fmt.Errorf("converter tarefas para JSON: %w", err)

	}

	errWrite := os.WriteFile(caminho, tarefasJSon, 0644)

	if errWrite != nil {
		return fmt.Errorf("Erro ao escrever no ficheiro: %w", errWrite)
	}
	return nil
}

// =====================================================================
// Exercício 2 — Carregar (os.ReadFile + json.Unmarshal)
// Carregar lê o ficheiro e devolve as tarefas.
//   - Se o ficheiro NÃO existir: devolve lista vazia e erro nil
//     (é a primeira execução, não é uma falha).
//   - Se existir mas o JSON for inválido: devolve erro não-nil.
//   - Qualquer outro erro de leitura: devolve-o embrulhado com %w.
//
// Dica: errors.Is(err, os.ErrNotExist)
// =====================================================================
func Carregar(caminho string) ([]Tarefa, error) {

	content, err := os.ReadFile(caminho)

	if err != nil {

		if errors.Is(err, os.ErrNotExist) {
			return []Tarefa{}, nil
		}

		return nil, fmt.Errorf("ler %s: %w", caminho, err)
	}

	var tarefas []Tarefa
	errObject := json.Unmarshal(content, &tarefas)
	if errObject != nil {

		return nil, fmt.Errorf("Error ao converter json to object: %w", errObject)
	}

	return tarefas, nil
}

// =====================================================================
// Exercício 3 — Gravar (os.Create + defer Close)
// Gravar escreve texto em caminho (substitui o que lá estiver).
// O erro do Close() também conta: se a escrita correu bem mas o Close
// falhar, esse erro tem de ser devolvido.
// Dica: resultado com nome (err error) + defer com closure, como no
// exercício 3/4 de exercicios4.
// =====================================================================
func Gravar(caminho, texto string) (err error) {

	f, err := os.Create(caminho)

	if err != nil {

		return fmt.Errorf("Error ao escrever no ficheiro %w", err)
	}

	defer func() {

		errClose := f.Close()

		if err == nil {

			err = errClose
		}

	}()

	_, err = f.WriteString(texto)

	if err != nil {

		return fmt.Errorf("Error al cerrar %w", err)
	}

	return nil
}

// =====================================================================
// Exercício 4 — Acrescentar (os.OpenFile + O_APPEND)
// Acrescentar adiciona linha + "\n" ao FIM do ficheiro, sem apagar o
// que já existe. Se o ficheiro não existir, cria-o.
// Flags: os.O_APPEND | os.O_CREATE | os.O_WRONLY, permissões 0644.
// =====================================================================
func Acrescentar(caminho, linha string) error {

	f, err := os.OpenFile(caminho, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if os.IsNotExist(err) {

		f, err = os.Create(caminho)

	}

	if err != nil {

		return fmt.Errorf("Erro ao abrir ficheiro para escrer %w", err)
	}

	defer func() {

		f.Close()

	}()

	_, err = f.WriteString(linha + "\n")

	if err != nil {

		return fmt.Errorf("Erro ao escrever no  ficheiro para  %w", err)

	}

	return nil
}

// =====================================================================
// Exercício 5 — LerLinhas (os.Open + bufio.Scanner)
// LerLinhas devolve as linhas do ficheiro (sem o "\n").
//   - Ficheiro vazio -> slice vazia, erro nil.
//   - Ficheiro inexistente -> erro tal que errors.Is(err, os.ErrNotExist)
//     seja true (se embrulhares, usa %w).
//
// Não esqueças o defer f.Close() logo após verificar o erro do Open, e
// verifica sc.Err() depois do ciclo.
// =====================================================================
func LerLinhas(caminho string) ([]string, error) {

	f, err := os.Open(caminho)

	if err != nil {

		return nil, fmt.Errorf("Erro ao validar o ficheiro %w", err)
	}

	defer func() {
		f.Close()
	}()

	sc := bufio.NewScanner(f)

	sliceString := make([]string, 0, 10)

	for sc.Scan() {

		sliceString = append(sliceString, sc.Text())
	}

	if err := sc.Err(); err != nil {

		return nil, fmt.Errorf("Erro ao validar o ficheiro %w", err)
	}

	return sliceString, nil
}

func main() {}
