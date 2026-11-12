package main

import (
	"log"
	"net/http"
)

func main() {
	// TODO (Etapa 6): escolher o repositório, criar o servidor e arrancar.
	s := NovoServidor(NovaMemoria())
	log.Println("a ouvir em :8080")
	log.Fatal(http.ListenAndServe(":8080", s.Handler()))
}
