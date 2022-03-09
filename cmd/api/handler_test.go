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

	type testCasesBlog struct {
		Body     string `json:"body,omitempty"`
		Category string `json:"category,omitempty"`
		Title    string `json:"title,omitempty"`
	}

	reqBody := data.Blog{
		Title:    "gRPC in Go!",
		Body:     "I do not know yet",
		Category: []string{"Golang", "Network"},
	}

	wantUrl := "/v1/blogs"

	tests := []struct {
		name         string
		urlPath      string
		wantCode     int
		reqBody      data.Blog
		envelopeName string
		wantBody     interface{}
	}{
		{name: "Must Success", urlPath: wantUrl, wantCode: http.StatusCreated,
			reqBody:      reqBody,
			envelopeName: "blog",
			wantBody: data.Blog{
				ID:        11,
				CreatedAt: time.Now(),
				Title:     "gRPC in Go!",
				Body:      "I do not know yet",
				Category:  []string{"Golang", "Network"},
				Version:   3,
			}},
		{name: "Empty Request", urlPath: wantUrl, wantCode: http.StatusUnprocessableEntity,
			envelopeName: "error",
			reqBody:      data.Blog{},
			wantBody: testCasesBlog{
				Body:     "must be provided",
				Title:    "must be provided",
				Category: "must be provided",
			}},
		{name: "Long Title", urlPath: wantUrl, wantCode: http.StatusUnprocessableEntity,
			envelopeName: "error",
			reqBody: data.Blog{
				Title:    "How to handle panics gracefully in Golang, How to handle panics gracefully in Golang, ",
				Body:     "I do not know yet",
				Category: []string{"Golang", "Network"},
			},
			wantBody: testCasesBlog{
				Title: "must not be more than 80 bytes long",
			}},
		{name: "Long Body", urlPath: wantUrl, wantCode: http.StatusUnprocessableEntity,
			envelopeName: "error",
			reqBody: data.Blog{
				Title:    "gRPC in Go!",
				Body:     string(make([]byte, 100001)),
				Category: []string{"Golang", "Network"},
			},
			wantBody: testCasesBlog{
				Body: "must not be more than 100000 bytes long",
			}},
		{name: "Min Category Size", urlPath: wantUrl, wantCode: http.StatusUnprocessableEntity,
			envelopeName: "error",
			reqBody: data.Blog{
				Title:    "gRPC in Go!",
				Body:     string(make([]byte, 1)),
				Category: []string{},
			},
			wantBody: testCasesBlog{
				Category: "must contain at least 1 categories",
			},
		},
		{name: "Max Category Size", urlPath: wantUrl, wantCode: http.StatusUnprocessableEntity,
			envelopeName: "error",
			reqBody: data.Blog{
				Title:    "gRPC in Go!",
				Body:     string(make([]byte, 1)),
				Category: []string{"Golang", "Network", "Distributed Systems", "Book", "RPC", "Complexity"},
			},
			wantBody: testCasesBlog{
				Category: "must not contain more than 5 categories",
			},
		},
		{name: "Unique Category", urlPath: wantUrl, wantCode: http.StatusUnprocessableEntity,
			envelopeName: "error",
			reqBody: data.Blog{
				Title:    "gRPC in Go!",
				Body:     string(make([]byte, 1)),
				Category: []string{"Golang", "Golang"},
			},
			wantBody: testCasesBlog{
				Category: "must not contain duplicate categories",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqBodyJSON, err := json.Marshal(tt.reqBody)
			if err != nil {
				t.Fatal(err)
			}

			wantBody, err := app.prettyJSON(envelope{tt.envelopeName: tt.wantBody})
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

			if !reflect.DeepEqual(wantBody, responseBody) {
				t.Errorf("Response Body -> want: \n%q; got: \n%q", wantBody, responseBody)
			}
		})
	}
}
