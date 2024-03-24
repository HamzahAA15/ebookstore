package customer

import (
	"context"
	"ebookstore/utils/request"
)

type ICustomerService interface {
	Register(ctx context.Context, customer *request.Register) (string, error)
	Login(ctx context.Context, customer request.Login) (string, error)
}
