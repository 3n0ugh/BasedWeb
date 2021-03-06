package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

// TODO: Add token model

type Model struct {
	Blog interface {
		Insert(blog *Blog) error
		Get(id int64) (*Blog, error)
		Update(blog *Blog) error
		Delete(id int64) error
		GetAll(title string, category []string, f Filter) ([]*Blog, Metadata, error)
	}
	User interface {
	}
}

func NewModel(db *sql.DB) Model {
	return Model{
		Blog: BlogModel{DB: db},
		User: UserModel{DB: db},
	}
}
