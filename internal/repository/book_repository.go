package repository

import (
	"ProjetAPI-GO-CRUD/internal/model"
	"context"
)

// BookRepository defines the contract for managing books in the database

type BookRepository interface {
	Create(ctx context.Context, book *model.CreateBookRequest) (int, error)
	GetByID(ctx context.Context, id int) (*model.Book, error)
	List(ctx context.Context) ([]model.Book, error)
	Update(ctx context.Context, id int, book *model.UpdateBookRequest) error
	Delete(ctx context.Context, id int) error
}
