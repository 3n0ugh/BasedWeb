package main

import (
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

type application struct{}

func main() {
	app := &application{}
	log.Fatal(http.ListenAndServe(":8080", app.routes()))
}

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("works"))
}

func (a *application) routes() http.Handler {
	router := httprouter.New()

	router.HandlerFunc(http.MethodGet, "/v1/health-check", HealthCheckHandler)
	return router
}
