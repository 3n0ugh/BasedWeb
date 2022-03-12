package main

import "net/http"

func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.logger.PrintInfo("request", map[string]string{
			"remoteAddress": r.RemoteAddr,
			"version":       r.Proto,
			"method":        r.Method,
			"requestURI":    r.URL.RequestURI()})

		next.ServeHTTP(w, r)
	})
}

// TODO: Add Rate Limitter
