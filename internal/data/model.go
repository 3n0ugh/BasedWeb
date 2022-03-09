package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
)

type Model struct {
	Blog interface {
		Insert(blog *Blog) error
		Get(id int64) (*Blog, error)
		Update(blog *Blog) error
		Delete(id int64) error
	}
}

func NewModel(db *sql.DB) Model {
	return Model{
		Blog: BlogModel{DB: db},
	}
}
