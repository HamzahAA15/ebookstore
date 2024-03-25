package service

import (
	"context"
	"ebookstore/internal/model/request"
	"ebookstore/internal/model/response"
)

type ICustomerService interface {
	Register(ctx context.Context, customer *request.Register) (string, error)
	Login(ctx context.Context, customer request.Login) (string, error)
}

type IOrderService interface {
	CreateOrder(ctx context.Context, req request.CreateOrder) (response.Order, error)
	GetUserOrders(ctx context.Context) ([]response.OrderData, error)
}

type IBookService interface {
	GetBooks(ctx context.Context) ([]response.Book, error)
}
