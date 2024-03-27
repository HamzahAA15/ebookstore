package customer_test

import (
	"context"
	"ebookstore/internal/model"
	"ebookstore/internal/model/request"
	"ebookstore/internal/repository/mocks"
	"ebookstore/internal/service/customer"
	mocksService "ebookstore/internal/service/mocks"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_customerService_Register(t *testing.T) {
	customerReq := request.Register{
		Email:    "email",
		Password: "password",
		Username: "username",
	}
	type fields struct {
		customerRepository  mocks.ICustomerRepository
		notificationService mocksService.INotificationService
	}
	type args struct {
		ctx      context.Context
		customer request.Register
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "best case",
			fields: fields{
				customerRepository: func() mocks.ICustomerRepository {
					m := mocks.ICustomerRepository{}
					m.On("GetCustomerByEmail", mock.Anything, customerReq.Email).Return(nil, nil)
					customerReq.Password = mock.Anything
					m.On("Register", mock.Anything, mock.Anything).Return(uint(1), nil)
					return m
				}(),
				notificationService: func() mocksService.INotificationService {
					m := mocksService.INotificationService{}
					m.On("SendEmail", mock.Anything, mock.Anything).Return(nil)
					return m
				}(),
			},
			args: args{
				ctx:      context.Background(),
				customer: customerReq,
			},
			wantErr: false,
		},
		{
			name: "GetCustomerByEmail error",
			fields: fields{
				customerRepository: func() mocks.ICustomerRepository {
					m := mocks.ICustomerRepository{}
					m.On("GetCustomerByEmail", mock.Anything, customerReq.Email).Return(nil, errors.New("error"))
					return m
				}(),
			},
			args: args{
				ctx:      context.Background(),
				customer: customerReq,
			},
			wantErr: true,
		},
		{
			name: "email already exist",
			fields: fields{
				customerRepository: func() mocks.ICustomerRepository {
					m := mocks.ICustomerRepository{}
					m.On("GetCustomerByEmail", mock.Anything, customerReq.Email).Return(&model.Customer{
						Email: customerReq.Email,
					}, nil)
					return m
				}(),
			},
			args: args{
				ctx:      context.Background(),
				customer: customerReq,
			},
			wantErr: true,
		},
		{
			name: "Register error",
			fields: fields{
				customerRepository: func() mocks.ICustomerRepository {
					m := mocks.ICustomerRepository{}
					m.On("GetCustomerByEmail", mock.Anything, customerReq.Email).Return(nil, nil)
					m.On("Register", mock.Anything, mock.Anything).Return(uint(0), errors.New("error"))
					return m
				}(),
			},
			args: args{
				ctx:      context.Background(),
				customer: customerReq,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := customer.NewCustomerService(&tt.fields.customerRepository, &tt.fields.notificationService)
			_, err := s.Register(tt.args.ctx, tt.args.customer)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}

func Test_customerService_Login(t *testing.T) {
	type fields struct {
		customerRepository  mocks.ICustomerRepository
		notificationService mocksService.INotificationService
	}
	type args struct {
		ctx      context.Context
		customer request.Login
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "best case",
			fields: fields{
				customerRepository: func() mocks.ICustomerRepository {
					m := mocks.ICustomerRepository{}
					m.On("GetCustomerByEmail", mock.Anything, mock.Anything).Return(&model.Customer{
						Email:    "email",
						Password: "$2a$12$KptVrUIFh4qX5.b8fHNjK.n1U749q8q86DtGxUFbEwbSUymQ./zty",
					}, nil)
					return m
				}(),
				notificationService: func() mocksService.INotificationService {
					m := mocksService.INotificationService{}
					m.On("SendEmail", mock.Anything, mock.Anything).Return(nil)
					return m
				}(),
			},
			args: args{
				ctx: context.Background(),
				customer: request.Login{
					Email:    "email",
					Password: "password",
				},
			},
			wantErr: false,
		},
		{
			name: "GetCustomerByEmail error",
			fields: fields{
				customerRepository: func() mocks.ICustomerRepository {
					m := mocks.ICustomerRepository{}
					m.On("GetCustomerByEmail", mock.Anything, mock.Anything).Return(nil, errors.New("error"))
					return m
				}(),
			},
			args: args{
				ctx: context.Background(),
				customer: request.Login{
					Email:    "email",
					Password: "password",
				},
			},
			wantErr: true,
		},
		{
			name: "invalid password",
			fields: fields{
				customerRepository: func() mocks.ICustomerRepository {
					m := mocks.ICustomerRepository{}
					m.On("GetCustomerByEmail", mock.Anything, mock.Anything).Return(&model.Customer{
						Email:    "email",
						Password: "password",
					}, nil)
					return m
				}(),
			},
			args: args{
				ctx: context.Background(),
				customer: request.Login{
					Email:    "email",
					Password: "password",
				},
			},
			wantErr: true,
		},
		{
			name: "invalid email",
			fields: fields{
				customerRepository: func() mocks.ICustomerRepository {
					m := mocks.ICustomerRepository{}
					m.On("GetCustomerByEmail", mock.Anything, mock.Anything).Return(&model.Customer{}, nil)
					return m
				}(),
			},
			args: args{
				ctx: context.Background(),
				customer: request.Login{
					Email:    "email",
					Password: "password",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := customer.NewCustomerService(&tt.fields.customerRepository, &tt.fields.notificationService)
			_, err := s.Login(tt.args.ctx, tt.args.customer)
			if (err != nil) != tt.wantErr {
				t.Errorf("customerService.Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
