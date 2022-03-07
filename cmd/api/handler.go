package main

import (
	"fmt"
	"github.com/3n0ugh/BasedWeb/internal/data"
	"net/http"
	"time"
)

func (app *application) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	err := app.writeJSON(w, http.StatusOK, envelope{"message": "works"}, r.Header)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	blog := &data.Blog{
		Title:    input.Title,
		Body:     input.Body,
		Category: input.Category,
	}

	err = app.model.Blog.Insert(blog)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/blogs/%d", blog.ID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"blog": blog}, headers)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
