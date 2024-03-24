package httpservice

import (
	"github.com/gofiber/fiber/v2"
)

func (h *BookHandler) SetupRoutes(app *fiber.App, auth fiber.Handler) {
	bookGroup := app.Group("/api/book")
	bookGroup.Get("/", auth, h.GetBooks)
}
