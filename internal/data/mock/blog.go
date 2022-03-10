package mock

import (
	"github.com/3n0ugh/BasedWeb/internal/data"
	"time"
)

var MockBlog = &data.Blog{
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
	blog.ID = MockBlog.ID
	blog.Version = MockBlog.Version
	blog.CreatedAt = MockBlog.CreatedAt
	return nil
}

func (b BlogModel) Get(id int64) (*data.Blog, error) {
	if id == MockBlog.ID {
		return MockBlog, nil
	}
	return nil, data.ErrRecordNotFound
}

func (b BlogModel) Update(blog *data.Blog) error {
	return nil
}

func (b BlogModel) Delete(id int64) error {
	if id == MockBlog.ID {
		return nil
	}
	return data.ErrRecordNotFound
}
