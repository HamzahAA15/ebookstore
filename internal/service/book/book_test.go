package book_test

import (
	"context"
	"ebookstore/internal/model"
	"ebookstore/internal/model/response"
	"ebookstore/internal/repository/mocks"
	"ebookstore/internal/service/book"
	"errors"
	"reflect"
	"testing"

	"github.com/stretchr/testify/mock"
)

func Test_bookService_GetBooks(t *testing.T) {
	type fields struct {
		bookRepository mocks.IBookRepository
	}

	type args struct {
		ctx context.Context
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []response.Book
		wantErr bool
	}{
		{
			name: "best case",
			fields: fields{
				bookRepository: func() mocks.IBookRepository {
					m := mocks.IBookRepository{}
					m.On("GetBooks", mock.Anything).Return([]model.Book{
						{
							ID:         uint(1),
							Title:      "title",
							Author:     "author",
							Price:      100,
							CategoryID: uint(1),
						},
					}, nil)
					m.On("GetCategoryByID", mock.Anything, uint(1)).Return(model.Category{
						ID:   1,
						Name: "category",
					}, nil)

					return m
				}(),
			},
			want: []response.Book{
				{
					ID:       1,
					Title:    "title",
					Author:   "author",
					Price:    100,
					Category: "category",
				},
			},
			wantErr: false,
		},
		{
			name: "failed GetBooks",
			fields: fields{
				bookRepository: func() mocks.IBookRepository {
					m := mocks.IBookRepository{}
					m.On("GetBooks", mock.Anything).Return([]model.Book{}, errors.New("failed"))
					return m
				}(),
			},
			want:    []response.Book{},
			wantErr: true,
		},
		{
			name: "failed GetCategoryByID",
			fields: fields{
				bookRepository: func() mocks.IBookRepository {
					m := mocks.IBookRepository{}
					m.On("GetBooks", mock.Anything).Return([]model.Book{
						{
							ID:         uint(1),
							Title:      "title",
							Author:     "author",
							Price:      100,
							CategoryID: uint(1),
						},
					}, nil)
					m.On("GetCategoryByID", mock.Anything, uint(1)).Return(model.Category{}, errors.New("failed"))
					return m
				}(),
			},
			want:    []response.Book{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := book.NewBookService(&tt.fields.bookRepository)
			got, err := s.GetBooks(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("bookService.GetBooks() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("bookService.GetBooks() = %v, want %v", got, tt.want)
			}
		})
	}
}
