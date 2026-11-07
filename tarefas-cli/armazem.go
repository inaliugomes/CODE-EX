package main

import (
	"encoding/json"
	"errors"
	"os"
)

// TODO (Etapa 2): Carregar(caminho string) (Lista, error)
//   - ficheiro inexistente -> lista vazia, sem erro
// TODO (Etapa 2): Guardar(caminho string, l Lista) error

func Carregar(caminho string) (Lista, error) {

	f, err := os.ReadFile(caminho)

	if errors.Is(err, os.ErrNotExist) {

		return Lista{}, nil
	}

	if err != nil {

		return Lista{}, err
	}

	var list Lista

	JsonObjeto := json.Unmarshal(f, &list)

	if JsonObjeto != nil {

		return Lista{}, JsonObjeto

	}

	return list, nil

}

func Guardar(caminho string, l Lista) error {

	jsonObject, err := json.MarshalIndent(l, "", " ")

	if err != nil {

		return err
	}

	errWriteJSon := os.WriteFile(caminho, jsonObject, 0644)

	if errWriteJSon != nil {
		return errWriteJSon
	}

	return nil

}
