package httpservice

import "github.com/gofiber/fiber/v2"

func (h *CustomerHandler) SetupRoutes(app *fiber.App) {
	orderGroup := app.Group("/api/customer")
	orderGroup.Post("/register", h.Register)
	orderGroup.Post("/login", h.Login)
}
