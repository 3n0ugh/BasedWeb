package main

import (
	"context"
	"encoding/json"
	"github.com/3n0ugh/BasedWeb/internal/data"
	"github.com/3n0ugh/BasedWeb/internal/data/mock"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"
)

type TestCases struct {
	name         string
	urlPath      string
	param        string
	wantCode     int
	reqBody      testBlog
	envelopeName string
	wantBody     interface{}
}

type testBlog struct {
	Title    string
	Body     string
	Category []string
}

type testCasesBlog struct {
	Body     string `json:"body,omitempty"`
	Category string `json:"category,omitempty"`
	Title    string `json:"title,omitempty"`
}

func NewTestApplication(model data.Model) *application {
	return &application{
		model: model,
	}
}

func NewRequestWithContext(method string, url string, bodyJSON []byte, p httprouter.Params) *http.Request {
	body := strings.NewReader(string(bodyJSON))

	r := httptest.NewRequest(method, url, body)

	ctx := r.Context()
	ctx = context.WithValue(ctx, httprouter.ParamsKey, p)
	r = r.WithContext(ctx)

	return r
}

func Check(t *testing.T, w *httptest.ResponseRecorder, tt TestCases) {
	body, err := ioutil.ReadAll(w.Result().Body)
	if err != nil {
		t.Fatal(err)
	}

	if tt.wantCode != w.Result().StatusCode {
		t.Errorf("Status Code -> want: %d; got: %d", tt.wantCode, w.Result().StatusCode)
	}

	if !reflect.DeepEqual(tt.wantBody, body) {
		t.Errorf("Body -> want: \n%q; got: \n%q", tt.wantBody, body)
	}
}

func TestHealthCheckHandler(t *testing.T) {
	app := &application{}

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/v1/health-check", nil)

	app.HealthCheckHandler(w, r)

	wantBody, err := app.prettyJSON(envelope{"message": "works"})
	if err != nil {
		t.Fatal(err)
	}

	test := TestCases{
		name:     "Health Check",
		urlPath:  "/v1/health-check",
		wantCode: http.StatusOK,
		wantBody: wantBody,
	}

	t.Run(test.name, func(t *testing.T) {
		Check(t, w, test)
	})
}

func TestCreateBlogHandler(t *testing.T) {
	app := NewTestApplication(mock.NewModel())

	reqBody := testBlog{
		Title:    "gRPC in Go!",
		Body:     "I do not know yet",
		Category: []string{"Golang", "Network"},
	}

	wantUrl := "/v1/blogs"

	tests := []TestCases{
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
			reqBody:      testBlog{},
			wantBody: testCasesBlog{
				Body:     "must be provided",
				Title:    "must be provided",
				Category: "must be provided",
			}},
		{name: "Long Title", urlPath: wantUrl, wantCode: http.StatusUnprocessableEntity,
			envelopeName: "error",
			reqBody: testBlog{
				Title:    "How to handle panics gracefully in Golang, How to handle panics gracefully in Golang, ",
				Body:     "I do not know yet",
				Category: []string{"Golang", "Network"},
			},
			wantBody: testCasesBlog{
				Title: "must not be more than 80 bytes long",
			}},
		{name: "Long Body", urlPath: wantUrl, wantCode: http.StatusUnprocessableEntity,
			envelopeName: "error",
			reqBody: testBlog{
				Title:    "gRPC in Go!",
				Body:     string(make([]byte, 100001)),
				Category: []string{"Golang", "Network"},
			},
			wantBody: testCasesBlog{
				Body: "must not be more than 100000 bytes long",
			}},
		{name: "Min Category Size", urlPath: wantUrl, wantCode: http.StatusUnprocessableEntity,
			envelopeName: "error",
			reqBody: testBlog{
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
			reqBody: testBlog{
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
			reqBody: testBlog{
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

			tt.wantBody, err = app.prettyJSON(envelope{tt.envelopeName: tt.wantBody})
			if err != nil {
				t.Fatal(err)
			}

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, tt.urlPath, strings.NewReader(string(reqBodyJSON)))

			app.createBlogHandler(w, r)

			Check(t, w, tt)
		})
	}
}

func TestShowBlogHandler(t *testing.T) {
	app := NewTestApplication(mock.NewModel())

	wantBodySuccess, err := app.prettyJSON(envelope{"blog": mock.Blog})
	if err != nil {
		t.Fatal(err)
	}

	tests := []TestCases{
		{
			name:    "Valid ID",
			urlPath: "/v1/blogs/",
			param:   "11", wantCode: http.StatusOK,
			wantBody: wantBodySuccess,
		},
		{
			name:     "Valid String ID",
			urlPath:  "/v1/blogs/\"11\"",
			param:    "11",
			wantCode: http.StatusOK,
			wantBody: wantBodySuccess,
		},
		{
			name:     "Negative ID",
			urlPath:  "/v1/blogs/-11",
			param:    "-11",
			wantCode: http.StatusBadRequest,
			wantBody: []byte("{\n\t\"error\": \"invalid id parameter\"\n}\n"),
		},
		{
			name:     "Non-existent ID",
			urlPath:  "/v1/blogs/12",
			param:    "12",
			wantCode: http.StatusNotFound,
			wantBody: []byte("{\n\t\"error\": \"the requested resource could not be found\"\n}\n"),
		},
		{
			name:     "Non-existent String ID",
			urlPath:  "/v1/blogs/\"12\"",
			param:    "12",
			wantCode: http.StatusNotFound,
			wantBody: []byte("{\n\t\"error\": \"the requested resource could not be found\"\n}\n"),
		},
		{
			name:     "Decimal ID",
			urlPath:  "/v1/blogs/1.11",
			param:    "1.11",
			wantCode: http.StatusBadRequest,
			wantBody: []byte("{\n\t\"error\": \"invalid id parameter\"\n}\n"),
		},
		{
			name:     "Empty ID",
			urlPath:  "/v1/blogs/",
			param:    "",
			wantCode: http.StatusBadRequest,
			wantBody: []byte("{\n\t\"error\": \"invalid id parameter\"\n}\n"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			w := httptest.NewRecorder()
			r := NewRequestWithContext(http.MethodGet, tt.urlPath, nil, httprouter.Params{
				{"id", tt.param},
			})

			app.showBlogHandler(w, r)

			Check(t, w, tt)
		})
	}
}

func TestDeleteBlogHandler(t *testing.T) {
	app := NewTestApplication(mock.NewModel())

	wantBodySuccess, err := app.prettyJSON(envelope{"message": "blogs successfully deleted"})
	if err != nil {
		t.Fatal(err)
	}

	tests := []TestCases{
		{
			name:    "Valid ID",
			urlPath: "/v1/blogs/",
			param:   "11", wantCode: http.StatusOK,
			wantBody: wantBodySuccess,
		},
		{
			name:     "Valid String ID",
			urlPath:  "/v1/blogs/\"11\"",
			param:    "11",
			wantCode: http.StatusOK,
			wantBody: wantBodySuccess,
		},
		{
			name:     "Negative ID",
			urlPath:  "/v1/blogs/-11",
			param:    "-11",
			wantCode: http.StatusBadRequest,
			wantBody: []byte("{\n\t\"error\": \"invalid id parameter\"\n}\n"),
		},
		{
			name:     "Non-existent ID",
			urlPath:  "/v1/blogs/12",
			param:    "12",
			wantCode: http.StatusNotFound,
			wantBody: []byte("{\n\t\"error\": \"the requested resource could not be found\"\n}\n"),
		},
		{
			name:     "Non-existent String ID",
			urlPath:  "/v1/blogs/\"12\"",
			param:    "12",
			wantCode: http.StatusNotFound,
			wantBody: []byte("{\n\t\"error\": \"the requested resource could not be found\"\n}\n"),
		},
		{
			name:     "Decimal ID",
			urlPath:  "/v1/blogs/1.11",
			param:    "1.11",
			wantCode: http.StatusBadRequest,
			wantBody: []byte("{\n\t\"error\": \"invalid id parameter\"\n}\n"),
		},
		{
			name:     "Empty ID",
			urlPath:  "/v1/blogs/",
			param:    "",
			wantCode: http.StatusBadRequest,
			wantBody: []byte("{\n\t\"error\": \"invalid id parameter\"\n}\n"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			w := httptest.NewRecorder()
			r := NewRequestWithContext(http.MethodDelete, tt.urlPath, nil, httprouter.Params{
				{"id", tt.param},
			})

			app.deleteBlogHandler(w, r)

			Check(t, w, tt)
		})
	}
}

func TestUpdateBlogHandler(t *testing.T) {
	app := NewTestApplication(mock.NewModel())

	wantBlog := mock.Blog
	wantBlog.Body = "I will learn."
	wantBlog.Title = "gRPC in Golang!"
	wantBlog.Category = []string{"Golang", "Network", "Framework"}
	wantBlog.Version = wantBlog.Version + 1

	tests := []TestCases{
		{
			name:     "Must Success",
			urlPath:  "/v1/blogs/11",
			param:    "11",
			wantCode: http.StatusOK,
			reqBody: testBlog{
				Title:    "gRPC in Golang!",
				Body:     "I will learn.",
				Category: []string{"Golang", "Network", "Framework"},
			},
			envelopeName: "blog",
			wantBody:     wantBlog,
		},
		{
			name:     "Negative ID",
			urlPath:  "/v1/blogs/-11",
			param:    "-11",
			wantCode: http.StatusBadRequest,
			reqBody: testBlog{
				Title:    "gRPC in Golang!",
				Body:     "I will learn.",
				Category: []string{"Golang", "Network", "Framework"},
			},
			envelopeName: "error",
			wantBody:     "invalid id parameter",
		},
		{
			name:     "Non-existent ID",
			urlPath:  "/v1/blogs/12",
			param:    "12",
			wantCode: http.StatusNotFound,
			reqBody: testBlog{
				Title:    "gRPC in Golang!",
				Body:     "I will learn.",
				Category: []string{"Golang", "Network", "Framework"},
			},
			envelopeName: "error",
			wantBody:     "the requested resource could not be found",
		},
		{
			name:     "Non-existent String ID",
			urlPath:  "/v1/blogs/\"12\"",
			param:    "12",
			wantCode: http.StatusNotFound,
			reqBody: testBlog{
				Title:    "gRPC in Golang!",
				Body:     "I will learn.",
				Category: []string{"Golang", "Network", "Framework"},
			},
			envelopeName: "error",
			wantBody:     "the requested resource could not be found",
		},
		{
			name:     "Decimal ID",
			urlPath:  "/v1/blogs/1.11",
			param:    "1.11",
			wantCode: http.StatusBadRequest,
			reqBody: testBlog{
				Title:    "gRPC in Golang!",
				Body:     "I will learn.",
				Category: []string{"Golang", "Network", "Framework"},
			},
			envelopeName: "error",
			wantBody:     "invalid id parameter",
		},
		{
			name:     "Empty ID",
			urlPath:  "/v1/blogs/",
			param:    "",
			wantCode: http.StatusBadRequest,
			reqBody: testBlog{
				Title:    "gRPC in Golang!",
				Body:     "I will learn.",
				Category: []string{"Golang", "Network", "Framework"},
			},
			envelopeName: "error",
			wantBody:     "invalid id parameter",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqBodyJSON, err := json.Marshal(tt.reqBody)
			if err != nil {
				t.Fatal(err)
			}

			tt.wantBody, err = app.prettyJSON(envelope{tt.envelopeName: tt.wantBody})
			if err != nil {
				t.Fatal(err)
			}

			w := httptest.NewRecorder()
			r := NewRequestWithContext(http.MethodPut, tt.urlPath, reqBodyJSON, httprouter.Params{
				{"id", tt.param},
			})

			app.updateBlogHandler(w, r)

			Check(t, w, tt)
		})
	}
}
