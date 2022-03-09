package main

import (
	"errors"
	"fmt"
	"github.com/3n0ugh/BasedWeb/internal/data"
	"github.com/3n0ugh/BasedWeb/internal/validator"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
	"time"
)

func (app *application) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	err := app.writeJSON(w, http.StatusOK, envelope{"message": "works"}, r.Header)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) createBlogHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		ID        int64     `json:"id"`
		CreatedAt time.Time `json:"-"`
		Title     string    `json:"title"`
		Body      string    `json:"body"`
		Category  []string  `json:"category"`
		Version   int32     `json:"version,omitempty"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	blog := &data.Blog{
		Title:    input.Title,
		Body:     input.Body,
		Category: input.Category,
	}

	v := validator.New()

	if data.ValidateBlog(v, blog); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.model.Blog.Insert(blog)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/blogs/%d", blog.ID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"blog": blog}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) showBlogHandler(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	if err != nil || id < 1 {
		app.badRequestResponse(w, r, err)
		return
	}

	blog, err := app.model.Blog.Get(id)
	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			app.notFoundResponse(w, r)
			return
		}
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"blog": blog}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

// TODO: Delete blog inorder to given id
func (app *application) deleteBlogHandler(w http.ResponseWriter, r *http.Request) {

}
