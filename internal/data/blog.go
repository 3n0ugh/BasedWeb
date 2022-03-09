package data

import (
	"context"
	"database/sql"
	"github.com/3n0ugh/BasedWeb/internal/validator"
	"github.com/lib/pq"
	"time"
)

type Blog struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"-"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	Category  []string  `json:"category"`
	Version   int32     `json:"version,omitempty"`
}

func ValidateBlog(v *validator.Validator, blog *Blog) {
	v.Check(blog.Title != "", "title", "must be provided")
	v.Check(len(blog.Title) <= 80, "title", "must not be more than 80 bytes long")

	v.Check(blog.Body != "", "body", "must be provided")
	v.Check(len(blog.Body) <= 1000000, "body", "must not be more than 100000 bytes long")

	v.Check(blog.Category != nil, "genres", "must be provided")
	v.Check(len(blog.Category) >= 1, "genres", "must contain at least 1 genre")
	v.Check(len(blog.Category) <= 5, "genres", "must not contain more than 5 genres")
	v.Check(validator.Unique(blog.Category), "genres", "must not contain duplicate genres")
}

type BlogModel struct {
	DB *sql.DB
}

func (b BlogModel) Insert(blog *Blog) error {
	query := `INSERT INTO blogs (title, body, category)
		VALUES ($1, $2, $3)
        RETURNING id, created_at, version`

	args := []interface{}{blog.Title, blog.Body, pq.Array(blog.Category)}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return b.DB.QueryRowContext(ctx, query, args...).
		Scan(&blog.ID, &blog.CreatedAt, &blog.Version)
}

func (b BlogModel) Get(id int64) (*Blog, error) {
	return nil, nil
}

func (b BlogModel) Update(blog *Blog) error {
	return nil
}

func (b BlogModel) Delete(id int64) error {
	return nil
}
