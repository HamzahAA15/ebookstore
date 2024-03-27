package customer_test

import (
	"bytes"
	"ebookstore/internal/httpservice/customer"
	"ebookstore/internal/model/request"
	"ebookstore/internal/model/response"
	"ebookstore/internal/service/mocks"
	"encoding/json"
	"errors"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCustomerHandler_Register(t *testing.T) {
	req := request.Register{
		Email:    "email@mail.com",
		Password: "Passw0rd.",
		Username: "username",
	}

	mockResp := response.Customer{}

	type fields struct {
		customerService mocks.ICustomerService
	}
	type args struct {
		c       *fiber.Ctx
		request request.Register
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
				customerService: func() mocks.ICustomerService {
					m := mocks.ICustomerService{}
					m.On("Register", mock.Anything, mock.Anything).Return("token", nil)
					return m
				}(),
			},
			args: args{
				c:       &fiber.Ctx{},
				request: req,
			},
			wantStatus: 200,
			wantMsg:    "success",
		},
		{
			name: "Register error",
			fields: fields{
				customerService: func() mocks.ICustomerService {
					m := mocks.ICustomerService{}
					m.On("Register", mock.Anything, mock.Anything).Return("", errors.New("error"))
					return m
				}(),
			},
			args: args{
				c:       &fiber.Ctx{},
				request: req,
			},
			wantStatus: 500,
			wantMsg:    "error",
		},
		{
			name: "invalid username",
			args: args{
				c: &fiber.Ctx{},
				request: request.Register{
					Username: "inv",
				},
			},
			wantStatus: 400,
			wantMsg:    "username",
		},
		{
			name: "invalid email",
			args: args{
				c: &fiber.Ctx{},
				request: request.Register{
					Username: "username",
					Email:    "invalid@mail",
				},
			},
			wantStatus: 400,
			wantMsg:    "email",
		},
		{
			name: "invalid password",
			args: args{
				c: &fiber.Ctx{},
				request: request.Register{
					Username: "username",
					Email:    "email@mail.com",
					Password: "invalid",
				},
			},
			wantStatus: 400,
			wantMsg:    "password",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := customer.NewCustomerHandler(&tt.fields.customerService)
			bodyBytes, _ := json.Marshal(tt.args.request)
			bodyIO := bytes.NewBuffer(bodyBytes)

			req := httptest.NewRequest("POST", "/customer/register", bodyIO)
			req.Header.Add("Content-Type", "application/json")
			srv := fiber.New()
			srv.Post("/customer/register", h.Register)

			resp, _ := srv.Test(req, 1000)

			bodyRespBytes, _ := io.ReadAll(resp.Body)
			json.Unmarshal(bodyRespBytes, &mockResp)

			assert.Contains(t, mockResp.Message, tt.wantMsg)
			assert.Equal(t, tt.wantStatus, resp.StatusCode)
		})
	}
}

func TestCustomerHandler_Login(t *testing.T) {
	type fields struct {
		customerService mocks.ICustomerService
	}
	type args struct {
		c   *fiber.Ctx
		req request.Login
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantStatus int
	}{
		{
			name: "best case",
			fields: fields{
				customerService: func() mocks.ICustomerService {
					m := mocks.ICustomerService{}
					m.On("Login", mock.Anything, mock.Anything).Return("token", nil)
					return m
				}(),
			},
			args: args{
				c:   &fiber.Ctx{},
				req: request.Login{Email: "username", Password: "Passw0rd."},
			},
			wantStatus: 200,
		},
		{
			name: "Login error",
			fields: fields{
				customerService: func() mocks.ICustomerService {
					m := mocks.ICustomerService{}
					m.On("Login", mock.Anything, mock.Anything).Return("", errors.New("error"))
					return m
				}(),
			},
			args: args{
				c:   &fiber.Ctx{},
				req: request.Login{Email: "username", Password: "Passw0rd."},
			},
			wantStatus: 500,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := customer.NewCustomerHandler(&tt.fields.customerService)
			bodyBytes, _ := json.Marshal(tt.args.req)
			bodyIO := bytes.NewBuffer(bodyBytes)

			req := httptest.NewRequest("POST", "/customer/login", bodyIO)
			req.Header.Add("Content-Type", "application/json")
			srv := fiber.New()
			srv.Post("/customer/login", h.Login)

			resp, _ := srv.Test(req, 1000)

			assert.Equal(t, tt.wantStatus, resp.StatusCode)
		})
	}
}
