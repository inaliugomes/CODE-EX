package main

import (
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	s := NovoServidor(NovoArmazem("tarefas.json"))

	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           s.Handler(),
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Println("a ouvir em :" + port)
	log.Fatal(srv.ListenAndServe()) // só devolve em erro
}
