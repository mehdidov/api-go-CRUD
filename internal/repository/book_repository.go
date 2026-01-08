package repository

import (
	"context"
	"database/sql"
	"fmt"

	"ProjetAPI-GO-CRUD/internal/model"
)

// BookRepository defines the contract for our database operations
type BookRepository interface {
	Create(ctx context.Context, book *model.CreateBookRequest) (int, error)
	GetByID(ctx context.Context, id int) (*model.Book, error)
	List(ctx context.Context) ([]model.Book, error)
	Update(ctx context.Context, id int, book *model.UpdateBookRequest) error
	Delete(ctx context.Context, id int) error
}

// postgresRepository is the concrete implementation using pure SQL
type postgresRepository struct {
	db *sql.DB
}

// NewPostgresRepository initializes a new instance of our repository
func NewPostgresRepository(db *sql.DB) BookRepository {
	return &postgresRepository{db: db}
}

// Create inserts a new book and returns the generated ID
func (r *postgresRepository) Create(ctx context.Context, req *model.CreateBookRequest) (int, error) {
	var newID int
	sqlQuery := `INSERT INTO books (title, author, year) VALUES ($1, $2, $3) RETURNING id`

	err := r.db.QueryRowContext(ctx, sqlQuery, req.Title, req.Author, req.Year).Scan(&newID)
	if err != nil {
		return 0, fmt.Errorf("failed to insert book: %w", err)
	}

	return newID, nil
}

// List fetches all books from the database
func (r *postgresRepository) List(ctx context.Context) ([]model.Book, error) {
	sqlQuery := `SELECT id, title, author, year FROM books`

	rows, err := r.db.QueryContext(ctx, sqlQuery)
	if err != nil {
		return nil, fmt.Errorf("could not fetch books list: %w", err)
	}
	defer rows.Close()

	var results []model.Book
	for rows.Next() {
		var b model.Book
		if err := rows.Scan(&b.ID, &b.Title, &b.Author, &b.Year); err != nil {
			return nil, err
		}
		results = append(results, b)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

// GetByID finds a single book by its primary key
func (r *postgresRepository) GetByID(ctx context.Context, id int) (*model.Book, error) {
	var b model.Book
	sqlQuery := `SELECT id, title, author, year FROM books WHERE id = $1`

	err := r.db.QueryRowContext(ctx, sqlQuery, id).Scan(&b.ID, &b.Title, &b.Author, &b.Year)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &b, nil
}

// Update modifies an existing book's data
func (r *postgresRepository) Update(ctx context.Context, id int, req *model.UpdateBookRequest) error {
	sqlQuery := `UPDATE books SET title = $1, author = $2, year = $3 WHERE id = $4`

	_, err := r.db.ExecContext(ctx, sqlQuery, req.Title, req.Author, req.Year, id)
	if err != nil {
		return fmt.Errorf("could not update book %d: %w", id, err)
	}

	return nil
}

// Delete removes a book entry from the database
func (r *postgresRepository) Delete(ctx context.Context, id int) error {
	sqlQuery := `DELETE FROM books WHERE id = $1`

	_, err := r.db.ExecContext(ctx, sqlQuery, id)
	if err != nil {
		return fmt.Errorf("failed to delete book %d: %w", id, err)
	}

	return nil
}
