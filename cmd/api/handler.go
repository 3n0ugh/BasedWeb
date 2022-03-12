package main

import (
	"errors"
	"fmt"
	"github.com/3n0ugh/BasedWeb/internal/data"
	"github.com/3n0ugh/BasedWeb/internal/validator"
	"net/http"
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
	id, err := app.readParamID(r)
	if err != nil {
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

func (app *application) deleteBlogHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readParamID(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	err = app.model.Blog.Delete(id)
	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			app.notFoundResponse(w, r)
			return
		}
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "blogs successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateBlogHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readParamID(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	blog, err := app.model.Blog.Get(id)
	if err != nil {
		if errors.Is(data.ErrRecordNotFound, err) {
			app.notFoundResponse(w, r)
			return
		}
		app.serverErrorResponse(w, r, err)
		return
	}

	var input struct {
		Title    *string  `json:"title"`
		Body     *string  `json:"body"`
		Category []string `json:"category"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.Title != nil {
		blog.Title = *input.Title
	}

	if input.Body != nil {
		blog.Body = *input.Body
	}

	if input.Category != nil {
		blog.Category = input.Category
	}

	v := validator.New()

	if data.ValidateBlog(v, blog); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.model.Blog.Update(blog)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"blog": blog}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) listBlogHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title    string
		Category []string
		data.Filter
	}

	v := validator.New()
	query := r.URL.Query()

	input.Title = app.readString(query, "title", "")

	input.Category = app.readCSV(query, "category", []string{})

	input.Filter.Page = app.readInt(query, "page", 1, v)
	input.Filter.PageSize = app.readInt(query, "page_size", 20, v)
	input.Filter.Sort = app.readString(query, "sort", "id")

	input.Filter.SortSafeList = []string{"id", "title", "-id", "-title"}

	if data.ValidateFilter(v, input.Filter); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	blogs, metadata, err := app.model.Blog.GetAll(input.Title, input.Category, input.Filter)
	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			app.notFoundResponse(w, r)
			return
		}
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"blogs": blogs, "metadata": metadata}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
