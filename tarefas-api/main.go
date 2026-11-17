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
	log.Println("a ouvir em :" + port)

	srv := &http.Server{Addr: ":" + port, Handler: s.Handler(), ReadHeaderTimeout: 5 * time.Second}

	err := srv.ListenAndServe()

	if err != nil {

		log.Fatal(err)
	}

}
