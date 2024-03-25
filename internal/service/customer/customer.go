package customer

import (
	"context"
	"ebookstore/internal/model"
	"ebookstore/internal/repository"
	authentication "ebookstore/utils/middleware"
	"ebookstore/utils/request"
	"errors"
	"fmt"
	"strings"
)

type customerService struct {
	customerRepository repository.ICustomerRepository
}

func NewCustomerService(customerRepository repository.ICustomerRepository) ICustomerService {
	return &customerService{
		customerRepository: customerRepository,
	}
}

func (s *customerService) Register(ctx context.Context, customer *request.Register) (string, error) {
	email := strings.ToLower(customer.Email)
	customerDB, err := s.customerRepository.GetCustomerByEmail(ctx, email)
	if err != nil {
		return "", fmt.Errorf("failed to get email existing: %s", err.Error())
	}

	if email == customerDB.Email {
		return "", errors.New("email already exists")
	}

	hashedPass, err := authentication.GenerateHashedPassword(customer.Password)
	if err != nil {
		return "", errors.New("failed to generate hashed password")
	}

	customerID, err := s.customerRepository.Register(ctx, &model.Customer{
		Email:    email,
		Password: hashedPass,
		Username: customer.Username,
	})

	if err != nil {
		return "", fmt.Errorf("failed to register customer: %s", err.Error())
	}

	token, err := authentication.GenerateToken(customer.Username, email, customerID)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %s", err.Error())
	}

	return token, nil
}

func (s *customerService) Login(ctx context.Context, customer request.Login) (string, error) {
	customerDB, err := s.customerRepository.GetCustomerByEmail(ctx, customer.Email)
	if err != nil {
		return "", fmt.Errorf("failed to login: %s", err.Error())
	}

	if customerDB.Email == "" {
		return "", errors.New("invalid email")
	}

	ok := authentication.CompareHashedPassword(customerDB.Password, customer.Password)
	if !ok {
		return "", errors.New("invalid password")
	}

	token, err := authentication.GenerateToken(customerDB.Username, customerDB.Email, customerDB.ID)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %s", err.Error())
	}

	return token, nil
}
