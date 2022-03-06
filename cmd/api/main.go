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

func (app *application) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	err := app.writeJSON(w, http.StatusOK, envelope{"message": "works"}, r.Header)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.HandlerFunc(http.MethodGet, "/v1/health-check", app.HealthCheckHandler)
	return router
}
