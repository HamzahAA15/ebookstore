package main

import (
	"ebookstore/internal/httpservice"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()
	httpservice.InitRoutes(app)

	app.Listen(":8080")
}
