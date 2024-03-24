package httpservice

import (
	"context"
	bookHandler "ebookstore/internal/httpservice/book"
	customerHandler "ebookstore/internal/httpservice/customer"
	orderHandler "ebookstore/internal/httpservice/order"
	"ebookstore/internal/repository"
	postgreRepo "ebookstore/internal/repository/postgresql"
	"ebookstore/internal/service/book"
	"ebookstore/internal/service/customer"
	"ebookstore/internal/service/order"
	authentication "ebookstore/utils/middleware"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "root"
	password = "password"
	dbname   = "book_db"
)

func InitRoutes(app *fiber.App) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err := repository.ConnectPostgres(context.Background(), psqlInfo)
	if err != nil {
		panic(err)
	}

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
