package main

import (
	"net/http"
	"reflect"
	"testing"
)

func TestHealthCheckHandler(t *testing.T) {
	app := newTestApplication(t)

	ts := newTestServer(t, app.routes())
	defer ts.Close()

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
		code, _, body := ts.get(t, test.urlPath)

		if code != test.wantCode {
			t.Errorf("status code -> want: %d; got: %d", test.wantCode, code)
		}
		if !reflect.DeepEqual(test.wantBody, body) {
			t.Errorf("body -> want: %q; got: %q", test.wantBody, body)
		}
	})
}
