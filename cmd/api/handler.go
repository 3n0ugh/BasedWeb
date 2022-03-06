package main

import "net/http"

func (app *application) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	err := app.writeJSON(w, http.StatusOK, envelope{"message": "works"}, r.Header)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
