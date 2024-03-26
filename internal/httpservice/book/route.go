package book

import (
	"github.com/gofiber/fiber/v2"
)

func (h *BookHandler) SetupRoutes(app *fiber.App) {
	bookGroup := app.Group("/api/book")
	bookGroup.Get("/", h.GetBooks)
}
