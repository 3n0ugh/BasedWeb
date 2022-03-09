package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.HandlerFunc(http.MethodGet, "/v1/health-check", app.HealthCheckHandler)

	router.HandlerFunc(http.MethodPost, "/v1/blogs", app.createBlogHandler)
	router.HandlerFunc(http.MethodGet, "/v1/blogs", app.showBlogHandler)

	return router
}
