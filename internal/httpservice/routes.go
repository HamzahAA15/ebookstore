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
	// gmailSMTP := gomail.NewDialer(config.CONFIG_SMTP_HOST, config.CONFIG_SMTP_PORT, config.CONFIG_AUTH_EMAIL, config.CONFIG_AUTH_PASSWORD)
	// notificationService := notification.NewGmailNotification(gmailSMTP)

	bookRepository := postgresql.NewBookRepository(db)
	bookService := bookService.NewBookService(bookRepository)
	bookHandler := bookHandler.NewBookHandler(bookService)
	bookHandler.SetupRoutes(app)

	customerRepository := postgresql.NewCustomerRepository(db)
	customerService := customerService.NewCustomerService(customerRepository)
	customerHandler := customerHandler.NewCustomerHandler(customerService)
	customerHandler.SetupRoutes(app)

	orderRepository := postgresql.NewOrderRepository(db)
	orderTxProvider := transactioner.NewTransactionProvider(db)
	orderService := orderService.NewOrderService(orderRepository, bookRepository, orderTxProvider)
	orderHandler := orderHandler.NewOrderHandler(orderService)
	orderHandler.SetupRoutes(app, auth)
}
