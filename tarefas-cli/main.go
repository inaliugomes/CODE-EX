package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func main() {
	// TODO (Etapa 3): chamar executar(os.Args[1:]) e imprimir o erro no stderr

	err := executar(os.Args[1:], "tarefas.json", os.Stdout)

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func executar(args []string, caminho string, saida io.Writer) error {

	if len(args) == 0 {

		fmt.Fprintln(saida, "uso: tarefas add|list|done|remove ...")
		return fmt.Errorf("texto vazio")
	}

	comando := strings.ToLower(args[0])

	listaTarefa, err := Carregar(caminho)

	if err != nil {

		return err
	}

	switch comando {

	case "add":

		if len(args) == 0 {

			fmt.Fprintln(saida, "uso: tarefas add|list|done|remove ...")
			return fmt.Errorf("texto vazio")
		}

		if len(args) < 2 {
			return fmt.Errorf("tarefa nao identificado correctamente")
		}
		tarefa, errTarefa := listaTarefa.Adicionar(args[1])

		if errTarefa != nil {
			return errTarefa
		}

		erroGuardar := Guardar(caminho, listaTarefa)

		if erroGuardar != nil {
			return erroGuardar
		}

		fmt.Fprintf(saida, "Tarefa %d criada.\n", tarefa.ID)

	case "remove":

		if len(args) == 0 {

			fmt.Fprintln(saida, "uso: tarefas add|list|done|remove ...")
			return fmt.Errorf("texto vazio")
		}

		if len(args) < 2 {
			return fmt.Errorf("tarefa nao identificado correctamente ")
		}

		idInt, errID := strconv.Atoi(args[1])

		if errID != nil {

			return errID
		}

		errorRemover := listaTarefa.Remover(idInt)

		if errorRemover != nil {

			return fmt.Errorf("erro ao remover a tarefa %w", errorRemover)
		}

		erroGuardar := Guardar(caminho, listaTarefa)

		if erroGuardar != nil {

			return erroGuardar

		}

		fmt.Fprintf(saida, "Tarefa %d removida.\n", idInt)

	case "done":

		if len(args) < 2 {

			return fmt.Errorf("tarefa nao identificado correctamente")

		}

		idInt, errID := strconv.Atoi(args[1])

		if errID != nil {

			return errID
		}

		errorConcluir := listaTarefa.Concluir(idInt)

		if errorConcluir != nil {
			return errorConcluir
		}

		erroGuardar := Guardar(caminho, listaTarefa)

		if erroGuardar != nil {

			return erroGuardar

		}

		fmt.Fprintf(saida, "Tarefa %d concluída.\n", idInt)

	case "list":

		for _, tarefa := range listaTarefa {

			var concluidoStr string

			if tarefa.Concluida {
				concluidoStr = "[x]"
			} else {

				concluidoStr = "[ ]"

			}

			fmt.Fprintf(saida, "%s %d  %s\n", concluidoStr, tarefa.ID, tarefa.Titulo)

		}

	default:
		return fmt.Errorf("comando nao valido para nenhuma operacao")

	}

	return nil

}
