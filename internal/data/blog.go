package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
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
	v.Check(len(blog.Body) <= 100000, "body", "must not be more than 100000 bytes long")

	v.Check(blog.Category != nil, "category", "must be provided")
	v.Check(len(blog.Category) >= 1, "category", "must contain at least 1 categories")
	v.Check(len(blog.Category) <= 5, "category", "must not contain more than 5 categories")
	v.Check(validator.Unique(blog.Category), "category", "must not contain duplicate categories")
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
	query := `SELECT created_at, title, body, category, version FROM blogs
		WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var blog Blog

	blog.ID = id

	row := b.DB.QueryRowContext(ctx, query, id)

	err := row.Scan(&blog.CreatedAt, &blog.Title, &blog.Body, pq.Array(&blog.Category), &blog.Version)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}

	return &blog, nil
}

func (b BlogModel) Update(blog *Blog) error {
	query := `UPDATE blogs
		SET title = $1, body = $2, category = $3, version = version + 1
		WHERE id = $4 AND version = $5
		RETURNING version`

	args := []interface{}{blog.Title, blog.Body, pq.Array(blog.Category), blog.ID, blog.Version}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := b.DB.QueryRowContext(ctx, query, args).Scan(&blog.Version)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrEditConflict
		}
		return err
	}
	return nil
}

func (b BlogModel) Delete(id int64) error {
	query := `DELETE FROM blogs
		WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	res, err := b.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	effectedRow, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if effectedRow == 0 {
		return ErrRecordNotFound
	}

	return nil
}

func (b BlogModel) GetAll(title string, category []string, f Filter) ([]*Blog, Metadata, error) {
	query := fmt.Sprintf(`
		SELECT count(*) OVER(), id, created_at, title, body, category, version
        FROM blogs
        WHERE (to_tsvector('simple', title) @@ plainto_tsquery('simple', $1)
		OR $1 = '')
        AND (category @> $2 OR $2 = '{}')
		ORDER BY %s %s, id ASC
		LIMIT $3 OFFSET $4`, f.sortColumn(), f.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	args := []interface{}{title, pq.Array(category), f.limit(), f.offset()}

	rows, err := b.DB.QueryContext(ctx, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, Metadata{}, ErrRecordNotFound
		}
		return nil, Metadata{}, err
	}

	defer func() {
		err = rows.Close()
		if err != nil {
			panic(err)
		}
	}()

	var totalRecords int
	var blogs = make([]*Blog, 1)

	for rows.Next() {
		var blog Blog

		err = rows.Scan(
			&totalRecords,
			&blog.ID,
			&blog.CreatedAt,
			&blog.Title,
			&blog.Body,
			pq.Array(&blog.Category),
			&blog.Version)

		if err != nil {
			return nil, Metadata{}, err
		}

		blogs = append(blogs, &blog)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, f.Page, f.PageSize)

	return blogs, metadata, nil
}
