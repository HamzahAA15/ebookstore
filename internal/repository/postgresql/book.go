package repository

import (
	"context"
	"ebookstore/internal/model"
	"ebookstore/internal/repository"

	"github.com/jmoiron/sqlx"
)

type BookRepository struct {
	db *sqlx.DB
}

func NewBookRepository(db *sqlx.DB) repository.IBookRepository {
	return &BookRepository{db: db}
}

func (r *BookRepository) GetBooks(ctx context.Context) ([]model.Book, error) {
	var books []model.Book
	query := "SELECT * FROM books"
	if err := r.db.SelectContext(ctx, &books, query); err != nil {
		return nil, err
	}

	return books, nil
}

func (r *BookRepository) GetBookByID(ctx context.Context, id uint) (model.Book, error) {
	var book model.Book
	query, err := r.db.PrepareContext(ctx, "SELECT id, title, author, price, category_id FROM books WHERE id = $1 AND deleted_at IS NULL")
	if err != nil {
		return model.Book{}, err
	}
	defer query.Close()

	query.QueryRowContext(ctx, id).Scan(&book.ID, &book.Title, &book.Author, &book.Price, &book.CategoryID)
	return book, nil
}
