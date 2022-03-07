package data

import "database/sql"

type Blog struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"-"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	Category  []string  `json:"category"`
	Version   int32     `json:"version,omitempty"`
}

type BlogModel struct {
	DB *sql.DB
}

func (b *BlogModel) Insert(blog *Blog) error {
	return nil
}

func (b *BlogModel) Get(id int64) (*Blog, error) {
	return nil, nil
}

func (b *BlogModel) Update(blog *Blog) error {
	return nil
}

func (b *BlogModel) Delete(id int64) error {
	return nil
}
