package book_test

import (
	"ebookstore/internal/httpservice/book"
	"ebookstore/internal/model/response"
	"ebookstore/internal/service/mocks"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestBookHandler_GetBooks(t *testing.T) {
	// fiberCtx := &fiber.Ctx{}
	type fields struct {
		bookService mocks.IBookService
	}
	type args struct {
		c *fiber.Ctx
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantStatus int
	}{
		{
			name: "best case",
			fields: fields{
				bookService: func() mocks.IBookService {
					m := mocks.IBookService{}
					m.On("GetBooks", mock.Anything).Return([]response.Book{
						{
							ID:       uint(1),
							Title:    "title",
							Author:   "author",
							Price:    100,
							Category: "category",
						},
					}, nil)
					return m
				}(),
			},
			args: args{
				c: &fiber.Ctx{},
			},
			wantStatus: 200,
		},
		{
			name: "GetBooks error",
			fields: fields{
				bookService: func() mocks.IBookService {
					m := mocks.IBookService{}
					m.On("GetBooks", mock.Anything).Return([]response.Book{}, errors.New("error"))
					return m
				}(),
			},
			args: args{
				c: &fiber.Ctx{},
			},
			wantStatus: 500,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := book.NewBookHandler(&tt.fields.bookService)
			req := httptest.NewRequest("GET", "/book", nil)
			req.Header.Add("Content-Type", "application/json")
			srv := fiber.New()
			srv.Get("/book", h.GetBooks)

			resp, err := srv.Test(req, 1000)
			if err != nil {
				println(err)
			}

			assert.Equal(t, tt.wantStatus, resp.StatusCode)
		})
	}
}
