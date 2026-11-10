package main

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
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

	caminhoDir := filepath.Dir(caminho)
	tmp, errPath := os.CreateTemp(caminhoDir, "tarefas-*.tmp")

	if errPath != nil {

		return errPath
	}

	_, errWriteJSon := tmp.Write(jsonObject)

	if errWriteJSon != nil {

		tmp.Close()
		os.Remove(tmp.Name())
		return errWriteJSon
	}
	errClose := tmp.Close()

	if errClose != nil {
		os.Remove(tmp.Name())
		return errClose
	}

	errRename := os.Rename(tmp.Name(), caminho)

	if errRename != nil {

		os.Remove(tmp.Name())
		return errRename
	}

	return nil

}
