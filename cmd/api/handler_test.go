package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestHealthCheckHandler(t *testing.T) {
	app := &application{}

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/v1/health-check", nil)

	app.HealthCheckHandler(w, r)

	wantBody, err := app.prettyJSON(envelope{"message": "works"})
	if err != nil {
		t.Fatal(err)
	}

	test := struct {
		name     string
		urlPath  string
		wantCode int
		wantBody []byte
	}{
		name:     "Health Check",
		urlPath:  "/v1/health-check",
		wantCode: http.StatusOK,
		wantBody: wantBody,
	}

	t.Run(test.name, func(t *testing.T) {

		body, err := ioutil.ReadAll(w.Body)
		if err != nil {
			t.Fatal(err)
		}

		if w.Code != test.wantCode {
			t.Errorf("status code -> want: %d; got: %d", test.wantCode, w.Code)
		}
		if !reflect.DeepEqual(test.wantBody, body) {
			t.Errorf("body -> want: %q; got: %q", test.wantBody, body)
		}
	})
}
