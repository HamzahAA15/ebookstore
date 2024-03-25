package httpservice

import (
	bookHandler "ebookstore/internal/httpservice/book"
	customerHandler "ebookstore/internal/httpservice/customer"
	orderHandler "ebookstore/internal/httpservice/order"
	bookService "ebookstore/internal/service/book"
	customerService "ebookstore/internal/service/customer"
	orderService "ebookstore/internal/service/order"

	"ebookstore/internal/repository/postgresql"
	authentication "ebookstore/utils/middleware"
	"ebookstore/utils/transactioner"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

func InitRoutes(app *fiber.App, db *sqlx.DB) {
	auth := authentication.AuthMiddleware()

	bookRepository := postgresql.NewBookRepository(db)
	bookService := bookService.NewBookService(bookRepository)
	bookHandler := bookHandler.NewBookHandler(bookService)
	bookHandler.SetupRoutes(app, auth)

	customerRepository := postgresql.NewCustomerRepository(db)
	customerService := customerService.NewCustomerService(customerRepository)
	customerHandler := customerHandler.NewCustomerHandler(customerService)
	customerHandler.SetupRoutes(app)

	orderRepository := postgresql.NewOrderRepository(db)
	orderTxProvider := transactioner.NewTransactionProvider(db)
	orderService := orderService.NewOrderService(orderRepository, bookRepository, orderTxProvider)
	orderHandler := orderHandler.NewOrderHandler(customerService, bookService, orderService)
	orderHandler.SetupRoutes(app, auth)
}
