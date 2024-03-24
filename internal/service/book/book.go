package book

import (
	"context"
	"ebookstore/internal/repository"
	"ebookstore/utils/response"
)

var catMap = map[uint]string{
	1:  "fiction",
	2:  "non-fiction",
	3:  "mystery",
	4:  "romance",
	5:  "science fiction",
	6:  "fantasy",
	7:  "biography",
	8:  "self-help",
	9:  "history",
	10: "thriller",
	11: "science",
}

type bookService struct {
	bookRepository repository.IBookRepository
}

func NewBookService(bookRepository repository.IBookRepository) IBookService {
	return &bookService{
		bookRepository: bookRepository,
	}
}

func (s *bookService) GetBooks(ctx context.Context) ([]response.Book, error) {
	resp := []response.Book{}
	books, err := s.bookRepository.GetBooks(ctx)
	if err != nil {
		return resp, err
	}

	for _, book := range books {
		resp = append(resp, response.Book{
			Title:    book.Title,
			Author:   book.Author,
			Price:    book.Price,
			Category: catMap[book.CategoryID],
		})
	}

	return resp, nil
}
