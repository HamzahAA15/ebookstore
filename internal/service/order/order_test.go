package order_test

import (
	"context"
	"ebookstore/internal/model"
	"ebookstore/internal/model/request"
	"ebookstore/internal/model/response"
	"ebookstore/internal/repository/mocks"
	mocksService "ebookstore/internal/service/mocks"
	"ebookstore/internal/service/order"
	"errors"
	"testing"
	"time"

	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_orderService_GetUserOrders(t *testing.T) {
	ctx := context.Background()
	id := uint(1)
	ctx = context.WithValue(ctx, "id", id)

	o := model.Order{
		ID:                1,
		CustomerID:        1,
		CustomerReference: "customerReference",
		ReceiverName:      "username",
		OrderDate:         time.Now().UTC().Truncate(time.Minute),
		Shipper:           "shipper",
		AirwaybillNumber:  "AWBnumber",
		TotalItem:         1,
		TotalPrice:        1,
		Address:           "address",
		City:              "city",
		District:          "district",
		PostalCode:        "postalCode",
	}

	item := model.Item{
		ID:       1,
		OrderID:  1,
		BookID:   1,
		Quantity: 1,
	}

	book := model.Book{
		ID:         1,
		Title:      "title",
		Author:     "author",
		Price:      10,
		CategoryID: 1,
	}

	orderData := []response.OrderData{
		{
			OrderID:           1,
			CustomerReference: "customerReference",
			ReceiverName:      "username",
			Address:           "address",
			City:              "city",
			District:          "district",
			PostalCode:        "postalCode",
			Shipper:           "shipper",
			AirwaybillNumber:  "AWBnumber",
			OrderDate:         time.Now().UTC().Truncate(time.Minute),
			TotalPrice:        1,
			TotalItem:         1,
		},
	}

	type fields struct {
		orderRepository     mocks.IOrderRepository
		TransactionProvider mocks.ITransactionProvider
		bookRepository      mocks.IBookRepository
		notificationService mocksService.INotificationService
	}

	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []response.OrderData
		wantErr bool
	}{
		{
			name: "best case",
			fields: fields{
				orderRepository: func() mocks.IOrderRepository {
					m := mocks.IOrderRepository{}
					m.On("GetOrderHistoryByCustomerID", mock.Anything, id).Return([]model.Order{o}, nil)
					m.On("GetItemsByOrderID", mock.Anything, o.ID).Return([]model.Item{item}, nil)
					orderData[0].Items = []response.Item{
						{
							BookID:   book.ID,
							Title:    book.Title,
							Author:   book.Author,
							Quantity: item.Quantity,
							Price:    book.Price,
						},
					}
					return m
				}(),
				bookRepository: func() mocks.IBookRepository {
					m := mocks.IBookRepository{}
					m.On("GetBookByID", mock.Anything, item.BookID).Return(book, nil)
					return m
				}(),
			},
			args: args{
				ctx: ctx,
			},
			want:    orderData,
			wantErr: false,
		},
		{
			name: "GetOrderHistoryByCustomerID error",
			fields: fields{
				orderRepository: func() mocks.IOrderRepository {
					m := mocks.IOrderRepository{}
					m.On("GetOrderHistoryByCustomerID", mock.Anything, id).Return([]model.Order{}, errors.New("error"))
					return m
				}(),
			},
			args: args{
				ctx: ctx,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "GetItemsByOrderID error",
			fields: fields{
				orderRepository: func() mocks.IOrderRepository {
					m := mocks.IOrderRepository{}
					m.On("GetOrderHistoryByCustomerID", mock.Anything, id).Return([]model.Order{o}, nil)
					m.On("GetItemsByOrderID", mock.Anything, o.ID).Return([]model.Item{}, errors.New("error"))
					return m
				}(),
			},
			args: args{
				ctx: ctx,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "GetBookByID error",
			fields: fields{
				orderRepository: func() mocks.IOrderRepository {
					m := mocks.IOrderRepository{}
					m.On("GetOrderHistoryByCustomerID", mock.Anything, id).Return([]model.Order{o}, nil)
					m.On("GetItemsByOrderID", mock.Anything, o.ID).Return([]model.Item{item}, nil)
					return m
				}(),
				bookRepository: func() mocks.IBookRepository {
					m := mocks.IBookRepository{}
					m.On("GetBookByID", mock.Anything, item.BookID).Return(model.Book{}, errors.New("error"))
					return m
				}(),
			},
			args: args{
				ctx: ctx,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := order.NewOrderService(&tt.fields.orderRepository, &tt.fields.bookRepository, &tt.fields.TransactionProvider, &tt.fields.notificationService)
			got, err := o.GetUserOrders(tt.args.ctx)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_orderService_CreateOrder(t *testing.T) {
	req := request.CreateOrder{
		Items: []request.Item{
			{
				BookID:   1,
				Quantity: 1,
			},
		},
		Address:      "address",
		City:         "city",
		District:     "district",
		PostalCode:   "postalCode",
		Shipper:      "shipper",
		ReceiverName: "username",
	}

	o := model.Order{
		ID:                1,
		CustomerID:        1,
		CustomerReference: mock.Anything,
		AirwaybillNumber:  mock.Anything,
		ReceiverName:      "username",
		OrderDate:         time.Now().UTC(),
		Shipper:           "shipper",
		Address:           "address",
		City:              "city",
		District:          "district",
		PostalCode:        "postalCode",
	}

	item := model.Item{
		OrderID:   1,
		BookID:    1,
		Quantity:  1,
		CreatedAt: time.Now().UTC().Truncate(time.Minute),
	}

	book := model.Book{
		ID:         1,
		Title:      "title",
		Author:     "author",
		Price:      10,
		CategoryID: 1,
	}

	ctx := context.Background()
	id := uint(1)
	ctx = context.WithValue(ctx, "id", id)
	ctx = context.WithValue(ctx, "email", "mail@mail.com")

	type fields struct {
		orderRepository     mocks.IOrderRepository
		TransactionProvider mocks.ITransactionProvider
		bookRepository      mocks.IBookRepository
		notificationService mocksService.INotificationService
	}

	type args struct {
		ctx context.Context
		req request.CreateOrder
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    response.CreateOrderData
		wantErr bool
	}{
		{
			name: "best case",
			fields: fields{
				TransactionProvider: func() mocks.ITransactionProvider {
					m := mocks.ITransactionProvider{}
					txProvide := mocks.TxxProvider{}
					txProvide.On("Commit").Return(nil)
					txProvide.On("Rollback").Return(nil)
					m.On("NewTransaction", mock.Anything).Return(&txProvide, nil)
					return m
				}(),
				orderRepository: func() mocks.IOrderRepository {
					m := mocks.IOrderRepository{}
					m.On("CreateOrder", mock.Anything, mock.Anything, mock.Anything).Return(o.ID, nil)
					o.ID = 1
					m.On("CreateItem", mock.Anything, mock.Anything, item).Return(nil)
					o.TotalItem = 1
					o.TotalPrice = 1
					o.UpdatedAt = pq.NullTime{Time: time.Now().UTC().Truncate(time.Minute), Valid: true}
					m.On("UpdateOrderByOrderID", mock.Anything, mock.Anything, mock.Anything).Return(nil)
					return m
				}(),
				bookRepository: func() mocks.IBookRepository {
					m := mocks.IBookRepository{}
					m.On("GetBookByID", mock.Anything, item.BookID).Return(book, nil)
					return m
				}(),
				notificationService: func() mocksService.INotificationService {
					m := mocksService.INotificationService{}
					m.On("SendNotification", mock.Anything, mock.Anything, mock.Anything).Return(nil)
					return m
				}(),
			},
			args: args{
				ctx: ctx,
				req: req,
			},
			want: response.CreateOrderData{
				OrderID:           o.ID,
				CustomerReference: o.CustomerReference,
				AirwaybillNumber:  "shipper",
				OrderDate:         o.OrderDate.String(),
			},
			wantErr: false,
		},
		{
			name: "NewTransaction error",
			fields: fields{
				TransactionProvider: func() mocks.ITransactionProvider {
					m := mocks.ITransactionProvider{}
					m.On("NewTransaction", mock.Anything).Return(nil, errors.New("error"))
					return m
				}(),
			},
			args: args{
				ctx: ctx,
				req: req,
			},
			want:    response.CreateOrderData{},
			wantErr: true,
		},
		{
			name: "CreateOrder error",
			fields: fields{
				TransactionProvider: func() mocks.ITransactionProvider {
					m := mocks.ITransactionProvider{}
					txProvide := mocks.TxxProvider{}
					txProvide.On("Commit").Return(nil)
					txProvide.On("Rollback").Return(nil)
					m.On("NewTransaction", mock.Anything).Return(&txProvide, nil)
					return m
				}(),
				orderRepository: func() mocks.IOrderRepository {
					m := mocks.IOrderRepository{}
					m.On("CreateOrder", mock.Anything, mock.Anything, mock.Anything).Return(uint(0), errors.New("error"))
					return m
				}(),
			},
			args: args{
				ctx: ctx,
				req: req,
			},
			want:    response.CreateOrderData{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := order.NewOrderService(&tt.fields.orderRepository, &tt.fields.bookRepository, &tt.fields.TransactionProvider, &tt.fields.notificationService)
			got, err := o.CreateOrder(tt.args.ctx, tt.args.req)
			if err != nil {
				println(err.Error())
			}
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Contains(t, got.AirwaybillNumber, tt.want.AirwaybillNumber)
		})
	}
}
