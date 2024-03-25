package postgresql

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
	query := "SELECT id, title, author, price, category_id FROM books WHERE deleted_at IS NULL"
	err := r.db.SelectContext(ctx, &books, query)
	if err != nil {
		return nil, err
	}

	return books, nil
}

func (r *BookRepository) GetBookByID(ctx context.Context, id uint) (model.Book, error) {
	var book model.Book
	query := ("SELECT id, title, author, price, category_id FROM books WHERE id = $1 AND deleted_at IS NULL")
	err := r.db.GetContext(ctx, &book, query, id)
	if err != nil {
		return model.Book{}, err
	}

	return book, nil
}

func (r *BookRepository) GetCategoryByID(ctx context.Context, id uint) (model.Category, error) {
	var category model.Category
	query := ("SELECT id, name FROM categories WHERE id = $1")
	err := r.db.GetContext(ctx, &category, query, id)
	if err != nil {
		return model.Category{}, err
	}

	return category, nil
}
