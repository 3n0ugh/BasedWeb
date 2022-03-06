package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

type config struct {
	port int
}

type application struct {
	config config
}

func main() {
	var cfg config
	var app application

	flag.IntVar(&cfg.port, "port", 8080, "API server port")
	flag.Parse()

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.port),
		Handler: app.routes(),
	}
	log.Fatal(srv.ListenAndServe())
}
