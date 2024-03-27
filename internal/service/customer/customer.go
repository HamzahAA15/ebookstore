package customer

import (
	"context"
	"ebookstore/internal/model"
	"ebookstore/internal/model/request"
	"ebookstore/internal/repository"
	"ebookstore/internal/service"
	"ebookstore/utils/config"
	authentication "ebookstore/utils/middleware"
	"ebookstore/utils/notification"
	"errors"
	"fmt"
	"log"
	"strings"
)

type customerService struct {
	customerRepository  repository.ICustomerRepository
	notificationService notification.INotificationService
}

func NewCustomerService(customerRepository repository.ICustomerRepository, notificationService notification.INotificationService) service.ICustomerService {
	return &customerService{
		customerRepository:  customerRepository,
		notificationService: notificationService,
	}
}

func (s *customerService) Register(ctx context.Context, customer request.Register) (string, error) {
	email := strings.ToLower(customer.Email)
	customerDB, err := s.customerRepository.GetCustomerByEmail(ctx, email)
	if err != nil {
		return "", fmt.Errorf("failed to get email existing: %s", err.Error())
	}

	if customerDB != nil && customerDB.Email == email {
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
	log.Println("qwe", config.CONFIG_EMAIL_SERVICE)
	//send notification
	if config.CONFIG_EMAIL_SERVICE {
		body := fmt.Sprintf(model.CustomerBodyEmailTemplate, customer.Email)

		emailPayload := notification.EmailPayload{
			To:      customer.Email,
			Subject: "Order Notification",
			Body:    body,
		}

		log.Printf("email %+v", emailPayload)

		go s.notificationService.SendNotification(emailPayload)
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
