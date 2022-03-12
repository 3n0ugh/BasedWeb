package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/v1/health-check", app.HealthCheckHandler)

	router.HandlerFunc(http.MethodPost, "/v1/blogs", app.createBlogHandler)
	router.HandlerFunc(http.MethodGet, "/v1/blogs/:id", app.showBlogHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/blogs/:id", app.deleteBlogHandler)
	router.HandlerFunc(http.MethodPut, "/v1/blogs/:id", app.updateBlogHandler)

	return app.logRequest(router)
}
