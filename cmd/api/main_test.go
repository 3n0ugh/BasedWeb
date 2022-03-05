package main

import (
	"bytes"
	"net/http"
	"testing"
)

func TestHealthCheckHandler(t *testing.T) {
	app := newTestApplication(t)
	router := app.routes()

	ts := newTestServer(t, router)
	defer ts.Close()

	test := struct {
		name     string
		urlPath  string
		wantCode int
		wantBody []byte
	}{
		name:     "Health Check",
		urlPath:  "/v1/health-check",
		wantCode: http.StatusOK,
		wantBody: []byte("works"),
	}

	t.Run(test.name, func(t *testing.T) {
		code, _, body := ts.get(t, test.urlPath)

		if code != test.wantCode {
			t.Errorf("status code -> want: %d; got: %d", test.wantCode, code)
		}
		if !bytes.Contains(body, test.wantBody) {
			t.Errorf("body -> want: %q; got: %q", test.wantBody, body)
		}
	})
}
