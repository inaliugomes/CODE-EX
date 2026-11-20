package main

import "errors"

var (
	ErrCodigoExiste = errors.New("código já existe")
	ErrNaoExiste    = errors.New("link não existe")
	ErrURLInvalido  = errors.New("url inválido")
)

type Link struct {
	Codigo  string `json:"codigo"`
	Destino string `json:"destino"`
	Visitas int    `json:"visitas"`
}
