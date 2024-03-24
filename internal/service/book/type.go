package book

import (
	"context"
	"ebookstore/utils/response"
)

type IBookService interface {
	GetBooks(ctx context.Context) ([]response.Book, error)
}
