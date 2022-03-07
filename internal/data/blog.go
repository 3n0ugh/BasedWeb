package data

import (
	"context"
	"database/sql"
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

type BlogModel struct {
	DB *sql.DB
}

func (b BlogModel) Insert(blog *Blog) error {
	query := `INSERT INTO blogs (title, body, category)
		VALUES ($1, $2)
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
