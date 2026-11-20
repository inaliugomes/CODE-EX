package main

import (
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	porta := os.Getenv("PORT")
	if porta == "" {
		porta = "8080"
	}

	servico := NovoServico(novaMemoria())

	srv := &http.Server{
		Addr:              ":" + porta,
		Handler:           NovoHandler(servico),
		ReadHeaderTimeout: 5 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}
