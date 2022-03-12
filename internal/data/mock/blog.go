package mock

import (
	"github.com/3n0ugh/BasedWeb/internal/data"
	"time"
)

var Blog = &data.Blog{
	ID:        11,
	CreatedAt: time.Now(),
	Title:     "gRPC in Go!",
	Body:      "I do not know yet",
	Category:  []string{"Golang", "Network"},
	Version:   3,
}

type BlogModel struct {
}

func (b BlogModel) Insert(blog *data.Blog) error {
	blog.ID = Blog.ID
	blog.Version = Blog.Version
	blog.CreatedAt = Blog.CreatedAt
	return nil
}

func (b BlogModel) Get(id int64) (*data.Blog, error) {
	if id == Blog.ID {
		return Blog, nil
	}
	return nil, data.ErrRecordNotFound
}

func (b BlogModel) Update(blog *data.Blog) error {
	if blog.ID != Blog.ID {
		return data.ErrEditConflict
	}
	return nil
}

func (b BlogModel) Delete(id int64) error {
	if id == Blog.ID {
		return nil
	}
	return data.ErrRecordNotFound
}

// TODO: Mock GetAll database function
