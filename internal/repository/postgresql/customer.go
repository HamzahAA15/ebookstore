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

func (c *customerRepository) Register(ctx context.Context, customer *model.Customer) error {
	statement, err := c.db.PrepareContext(ctx, "INSERT INTO customers (email, password, username) VALUES ($1, $2, $3)")
	if err != nil {
		return err
	}
	defer statement.Close()

	_, err = statement.ExecContext(ctx, customer.Email, customer.Password, customer.Username)
	if err != nil {
		return err
	}

	return nil
}

func (c *customerRepository) GetCustomerByEmail(ctx context.Context, email string) (*model.Customer, error) {
	var customer model.Customer
	query, err := c.db.PrepareContext(ctx, "SELECT id, email, username, password FROM customers WHERE email = $1")
	if err != nil {
		return nil, err
	}
	defer query.Close()

	err = query.QueryRowContext(ctx, email).Scan(&customer.ID, &customer.Email, &customer.Username, &customer.Password)
	if err != nil && !(err.Error() == "sql: no rows in result set") {
		return nil, err
	}

	return &customer, nil
}
