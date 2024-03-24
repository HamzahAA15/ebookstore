package order

import (
	"context"
	"ebookstore/utils/request"
	"ebookstore/utils/response"
)

type IOrderService interface {
	CreateOrder(ctx context.Context, req request.CreateOrder) (response.Order, error)
	GetUserOrders(ctx context.Context) ([]response.OrderData, error)
}
