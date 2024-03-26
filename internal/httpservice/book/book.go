package book

import (
	"ebookstore/internal/model/response"
	"ebookstore/internal/service"

	"github.com/gofiber/fiber/v2"
)

type BookHandler struct {
	bookService service.IBookService
}

func NewBookHandler(bookService service.IBookService) *BookHandler {
	return &BookHandler{
		bookService: bookService,
	}
}

func (h *BookHandler) GetBooks(c *fiber.Ctx) error {
	books, err := h.bookService.GetBooks(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.GetBooks{
			StatusCode: fiber.StatusInternalServerError,
			Message:    err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(response.GetBooks{
		StatusCode: fiber.StatusOK,
		Message:    "success",
		Data:       books,
	})
}
