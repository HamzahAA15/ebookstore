package repository

import (
	"context"
	"ebookstore/internal/model"
	"ebookstore/internal/repository"

	"github.com/jmoiron/sqlx"
)

type customerRepository struct {
	db *sqlx.DB
}

func NewCustomerRepository(db *sqlx.DB) repository.ICustomerRepository {
	return &customerRepository{db: db}
}

func (c *customerRepository) Register(ctx context.Context, customer *model.Customer) (uint, error) {
	var id uint
	query := "INSERT INTO customers (email, password, username) VALUES (:email, :password, :username) RETURNING id;"
	err := c.db.QueryRowxContext(ctx, query, customer).Scan(&id)
	if err != nil {
		return 0, err

	}

	return id, nil
}

func (c *customerRepository) GetCustomerByEmail(ctx context.Context, email string) (*model.Customer, error) {
	var customer model.Customer
	query := "SELECT id, email, username, password FROM customers WHERE email = $1"

	err := c.db.GetContext(ctx, &customer, query, email)
	if err != nil && !(err.Error() == "sql: no rows in result set") {
		return nil, err
	}

	return &customer, nil
}
