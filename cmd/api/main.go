package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/v1/health-check", HealthCheckHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("works"))
}
