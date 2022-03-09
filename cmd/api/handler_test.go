package main

import (
	"encoding/json"
	"github.com/3n0ugh/BasedWeb/internal/data"
	"github.com/3n0ugh/BasedWeb/internal/data/mock"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"
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

func TestCreateBlogHandler(t *testing.T) {
	app := &application{
		model: mock.NewModel(),
	}

	reqBody := data.Blog{
		Title:    "gRPC in Go!",
		Body:     "I do not know yet",
		Category: []string{"Golang", "Network"},
	}

	var blog = data.Blog{
		ID:        11,
		CreatedAt: time.Now(),
		Title:     "gRPC in Go!",
		Body:      "I do not know yet",
		Category:  []string{"Golang", "Network"},
		Version:   3,
	}

	wantBody, err := app.prettyJSON(envelope{"blog": blog})
	if err != nil {
		t.Fatal(err)
	}

	wantUrl := "/v1/blogs"

	tests := []struct {
		name     string
		urlPath  string
		wantCode int
		wantBody []byte
	}{
		// TODO: add some test cases after the write validation for blog
		{name: "Should Success", urlPath: wantUrl, wantCode: 201, wantBody: wantBody},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqBodyJSON, err := json.Marshal(reqBody)
			if err != nil {
				t.Fatal(err)
			}

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, tt.urlPath, strings.NewReader(string(reqBodyJSON)))

			app.createBlogHandler(w, r)

			responseBody, err := ioutil.ReadAll(w.Result().Body)
			if err != nil {
				t.Fatal(err)
			}

			if tt.wantCode != w.Result().StatusCode {
				t.Errorf("Status Code -> want: %d; got: %d", tt.wantCode, w.Result().StatusCode)
			}

			if !reflect.DeepEqual(tt.wantBody, responseBody) {
				t.Errorf("Response Body -> want: \n%q; got: \n%q", tt.wantBody, responseBody)
			}
		})
	}
}
