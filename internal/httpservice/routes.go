package httpservice

import (
	bookHandler "ebookstore/internal/httpservice/book"
	customerHandler "ebookstore/internal/httpservice/customer"
	orderHandler "ebookstore/internal/httpservice/order"
	postgreRepo "ebookstore/internal/repository/postgresql"
	"ebookstore/internal/service/book"
	"ebookstore/internal/service/customer"
	"ebookstore/internal/service/order"
	authentication "ebookstore/utils/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

func InitRoutes(app *fiber.App, db *sqlx.DB) {
	auth := authentication.AuthMiddleware()

	bookRepository := postgreRepo.NewBookRepository(db)
	bookService := book.NewBookService(bookRepository)
	bookHandler := bookHandler.NewBookHandler(bookService)
	bookHandler.SetupRoutes(app, auth)

	customerRepository := postgreRepo.NewCustomerRepository(db)
	customerService := customer.NewCustomerService(customerRepository)
	customerHandler := customerHandler.NewCustomerHandler(customerService)
	customerHandler.SetupRoutes(app)

	orderRepository := postgreRepo.NewOrderRepository(db)
	orderService := order.NewOrderService(orderRepository, bookRepository)
	orderHandler := orderHandler.NewOrderHandler(customerService, bookService, orderService)
	orderHandler.SetupRoutes(app, auth)
}
