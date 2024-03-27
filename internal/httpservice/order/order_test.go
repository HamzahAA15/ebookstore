package order_test

import (
	"bytes"
	"ebookstore/internal/httpservice/order"
	"ebookstore/internal/model/request"
	"ebookstore/internal/model/response"
	"ebookstore/internal/service"
	"ebookstore/internal/service/mocks"
	"encoding/json"
	"errors"
	"io"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestOrderHandler_CreateOrder(t *testing.T) {
	req := request.CreateOrder{
		ReceiverName: "hamzah",
		Address:      "address",
		City:         "city",
		District:     "district",
		PostalCode:   "postalCode",
		Shipper:      "shipper",
		Items: []request.Item{
			{
				BookID:   1,
				Quantity: 5,
			},
		},
	}

	data := response.CreateOrderData{
		OrderID:           1,
		CustomerReference: "customerReference",
		AirwaybillNumber:  "AWBNumber",
		OrderDate:         time.Now().UTC().Truncate(time.Minute).String(),
	}

	mockResp := response.Order{
		StatusCode: 200,
		Message:    "success",
		Data:       data,
	}

	type fields struct {
		orderService service.IOrderService
	}

	type args struct {
		request     request.CreateOrder
		authHandler fiber.Handler
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantStatus int
		wantMsg    string
	}{
		{
			name: "best case",
			fields: fields{
				orderService: func() *mocks.IOrderService {
					m := mocks.IOrderService{}
					m.On("CreateOrder", mock.Anything, req).Return(data, nil)
					return &m
				}(),
			},
			args: args{
				authHandler: func(c *fiber.Ctx) error {
					c.Locals("username", "hamzah")
					return c.Next()
				},
				request: req,
			},
			wantStatus: 200,
			wantMsg:    "success",
		},
		{
			name: "invalid request, items empty",
			args: args{
				authHandler: func(c *fiber.Ctx) error {
					c.Locals("username", "hamzah")
					return c.Next()
				},
				request: request.CreateOrder{
					Items: []request.Item{},
				},
			},
			wantStatus: 400,
			wantMsg:    "items",
		},
		{
			name: "invalid request, receiver address empty",
			args: args{
				authHandler: func(c *fiber.Ctx) error {
					c.Locals("username", "hamzah")
					return c.Next()
				},
				request: request.CreateOrder{
					Items: []request.Item{
						{
							BookID:   1,
							Quantity: 2,
						},
					},
				},
			},
			wantStatus: 400,
			wantMsg:    "address",
		},
		{
			name: "invalid request, receiver city empty",
			args: args{
				authHandler: func(c *fiber.Ctx) error {
					c.Locals("username", "hamzah")
					return c.Next()
				},
				request: request.CreateOrder{
					Items: []request.Item{
						{
							BookID:   1,
							Quantity: 2,
						},
					},
					ReceiverName: "hamzah",
					Address:      "address",
				},
			},
			wantStatus: 400,
			wantMsg:    "city",
		},
		{
			name: "invalid request, receiver district empty",
			args: args{
				authHandler: func(c *fiber.Ctx) error {
					c.Locals("username", "hamzah")
					return c.Next()
				},
				request: request.CreateOrder{
					Items: []request.Item{
						{
							BookID:   1,
							Quantity: 2,
						},
					},
					ReceiverName: "hamzah",
					Address:      "address",
					City:         "city",
				},
			},
			wantStatus: 400,
			wantMsg:    "district",
		},
		{
			name: "invalid request, receiver postal code empty",
			args: args{
				authHandler: func(c *fiber.Ctx) error {
					c.Locals("username", "hamzah")
					return c.Next()
				},
				request: request.CreateOrder{
					Items: []request.Item{
						{
							BookID:   1,
							Quantity: 2,
						},
					},
					ReceiverName: "hamzah",
					Address:      "address",
					City:         "city",
					District:     "district",
				},
			},
			wantStatus: 400,
			wantMsg:    "postal",
		},
		{
			name: "CreateOrder error",
			fields: fields{
				orderService: func() *mocks.IOrderService {
					m := mocks.IOrderService{}
					m.On("CreateOrder", mock.Anything, req).Return(data, errors.New("CreateOrder error"))
					return &m
				}(),
			},
			args: args{
				authHandler: func(c *fiber.Ctx) error {
					c.Locals("username", "hamzah")
					return c.Next()
				},
				request: req,
			},
			wantStatus: 500,
			wantMsg:    "error",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := order.NewOrderHandler(tt.fields.orderService)
			bodyBytes, _ := json.Marshal(tt.args.request)
			bodyIO := bytes.NewBuffer(bodyBytes)

			req := httptest.NewRequest("POST", "/order", bodyIO)
			req.Header.Add("Content-Type", "application/json")
			srv := fiber.New()
			srv.Post("/order", tt.args.authHandler, h.CreateOrder)

			resp, _ := srv.Test(req, 1000)
			bodyRespBytes, _ := io.ReadAll(resp.Body)
			json.Unmarshal(bodyRespBytes, &mockResp)

			assert.Contains(t, mockResp.Message, tt.wantMsg)
			assert.Equal(t, tt.wantStatus, resp.StatusCode)
		})
	}
}

func TestOrderHandler_GetUserOrders(t *testing.T) {
	type fields struct {
		orderService service.IOrderService
	}

	mockResp := response.Order{}

	data := []response.OrderData{
		{
			OrderID:           1,
			CustomerReference: "customerRef",
			ReceiverName:      "hamzah",
			Address:           "address",
			City:              "city",
			District:          "district",
			PostalCode:        "postalCode",
			Shipper:           "shipper",
			AirwaybillNumber:  "airwaybillNumber",
			OrderDate:         time.Now().UTC().Truncate(time.Minute),
			Items: []response.Item{
				{
					BookID:   1,
					Title:    "title",
					Author:   "author",
					Quantity: 5,
					Price:    10,
				},
			},
			TotalItem:  10,
			TotalPrice: 100,
		},
	}
	type args struct {
		authHandler fiber.Handler
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantStatus int
		wantMsg    string
	}{
		{
			name: "best case",
			fields: fields{
				orderService: func() *mocks.IOrderService {
					m := mocks.IOrderService{}
					m.On("GetUserOrders", mock.Anything, mock.Anything).Return(data, nil)
					return &m
				}(),
			},
			args: args{
				authHandler: func(c *fiber.Ctx) error {
					c.Locals("username", "hamzah")
					c.Locals("id", uint(1))
					return c.Next()
				},
			},
			wantStatus: 200,
			wantMsg:    "success",
		},
		{
			name: "GetUserOrders error",
			fields: fields{
				orderService: func() *mocks.IOrderService {
					m := mocks.IOrderService{}
					m.On("GetUserOrders", mock.Anything, mock.Anything).Return(data, errors.New("GetUserOrders error"))
					return &m
				}(),
			},
			args: args{
				authHandler: func(c *fiber.Ctx) error {
					c.Locals("username", "hamzah")
					c.Locals("id", uint(1))
					return c.Next()
				},
			},
			wantStatus: 500,
			wantMsg:    "error",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := order.NewOrderHandler(tt.fields.orderService)

			req := httptest.NewRequest("GET", "/order/order-history", nil)
			req.Header.Add("Content-Type", "application/json")
			srv := fiber.New()
			srv.Get("/order/order-history", tt.args.authHandler, h.GetUserOrders)

			resp, _ := srv.Test(req, 1000)
			bodyRespBytes, _ := io.ReadAll(resp.Body)
			json.Unmarshal(bodyRespBytes, &mockResp)

			assert.Contains(t, mockResp.Message, tt.wantMsg)
			assert.Equal(t, tt.wantStatus, resp.StatusCode)
		})
	}
}
