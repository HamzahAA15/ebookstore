package book

import (
	"context"
	"ebookstore/internal/model/response"
	"ebookstore/internal/repository"
	"ebookstore/internal/service"
	"fmt"
)

var catMap = make(map[uint]string)

type bookService struct {
	bookRepository repository.IBookRepository
}

func NewBookService(bookRepository repository.IBookRepository) service.IBookService {
	return &bookService{
		bookRepository: bookRepository,
	}
}

func (s *bookService) GetBooks(ctx context.Context) ([]response.Book, error) {
	resp := []response.Book{}
	books, err := s.bookRepository.GetBooks(ctx)
	if err != nil {
		return resp, fmt.Errorf("failed to get books: %s", err.Error())
	}

	for _, book := range books {
		var categoryName string
		//check inMemory caching
		categoryName, ok := catMap[book.CategoryID]
		if !ok {
			category, err := s.bookRepository.GetCategoryByID(ctx, book.CategoryID)
			if err != nil {
				return resp, fmt.Errorf("failed to get category: %s", err.Error())
			}

			categoryName = category.Name
		}

		resp = append(resp, response.Book{
			ID:       book.ID,
			Title:    book.Title,
			Author:   book.Author,
			Price:    book.Price,
			Category: categoryName,
		})
	}

	return resp, nil
}
