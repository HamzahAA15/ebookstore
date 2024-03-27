package repository

import (
	"context"
	"ebookstore/internal/model"
	"ebookstore/utils/transactioner"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type IBookRepository interface {
	GetBooks(ctx context.Context) ([]model.Book, error)
	GetBookByID(ctx context.Context, id uint) (model.Book, error)

	GetCategoryByID(ctx context.Context, id uint) (model.Category, error)
}

type ICustomerRepository interface {
	Register(ctx context.Context, customer *model.Customer) (uint, error)
	GetCustomerByEmail(ctx context.Context, email string) (*model.Customer, error)
}

type IOrderRepository interface {
	CreateOrder(ctx context.Context, tx transactioner.TxxProvider, order model.Order) (uint, error)
	GetOrderHistoryByCustomerID(ctx context.Context, customerID uint) ([]model.Order, error)
	UpdateOrderByOrderID(ctx context.Context, tx transactioner.TxxProvider, order model.Order) error

	CreateItem(ctx context.Context, tx transactioner.TxxProvider, item model.Item) error
	GetItemsByOrderID(ctx context.Context, orderID uint) ([]model.Item, error)
}

func ConnectPostgres(ctx context.Context, dbURL string) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Postgres: %w", err)
	}

	if err = db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping Postgres: %w", err)
	}

	log.Println("Connected to Postgres!")
	return db, nil
}
