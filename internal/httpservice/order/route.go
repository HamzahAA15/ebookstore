package order

import "github.com/gofiber/fiber/v2"

func (h *OrderHandler) SetupRoutes(app *fiber.App, auth fiber.Handler) {
	orderGroup := app.Group("/api/order")
	orderGroup.Post("/", auth, h.CreateOrder)
	orderGroup.Get("/order-history", auth, h.GetUserOrders)
}
