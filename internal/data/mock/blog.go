package mock

import (
	"github.com/3n0ugh/BasedWeb/internal/data"
	"time"
)

var mockBlog = &data.Blog{
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
	blog.ID = mockBlog.ID
	blog.Version = mockBlog.Version
	blog.CreatedAt = mockBlog.CreatedAt
	return nil
}

func (b BlogModel) Get(id int64) (*data.Blog, error) {
	return nil, nil
}

func (b BlogModel) Update(blog *data.Blog) error {
	return nil
}

func (b BlogModel) Delete(id int64) error {
	return nil
}
