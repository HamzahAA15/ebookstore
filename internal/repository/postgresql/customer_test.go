package postgresql_test

import (
	"context"
	"ebookstore/internal/model"
	"ebookstore/internal/repository/postgresql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func Test_customerRepository_Register(t *testing.T) {
	customerID := uint(1)
	type fields struct {
		data uint
		err  error
	}
	type args struct {
		ctx      context.Context
		customer *model.Customer
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    uint
		wantErr bool
	}{
		{
			name: "best case",
			fields: fields{
				data: customerID,
				err:  nil,
			},
			args: args{
				ctx:      context.Background(),
				customer: &model.Customer{},
			},
			want:    1,
			wantErr: false,
		},
		{
			name: "QueryRowxContext error",
			fields: fields{
				data: 0,
				err:  errors.New("some error"),
			},
			args: args{
				ctx:      context.Background(),
				customer: &model.Customer{},
			},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, m, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			sqlxDB := sqlx.NewDb(db, "sqlmock")
			testDB := postgresql.NewCustomerRepository(sqlxDB)

			query := "INSERT INTO customers (email, password, username) VALUES ($1, $2, $3) RETURNING id;"

			mockExpectQuery := m.ExpectQuery(query).WithArgs(tt.args.customer.Email, tt.args.customer.Password, tt.args.customer.Username)
			if tt.fields.err != nil {
				mockExpectQuery.WillReturnError(err)
			} else {
				row := sqlmock.NewRows([]string{"id"}).AddRow(tt.fields.data)
				mockExpectQuery.WillReturnRows(row)
			}

			got, err := testDB.Register(tt.args.ctx, tt.args.customer)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_customerRepository_GetCustomerByEmail(t *testing.T) {
	customer := model.Customer{
		ID:       1,
		Email:    "username@mail.com",
		Username: "usernamer",
		Password: "Passw0rd.",
	}

	type fields struct {
		data model.Customer
		err  error
	}
	type args struct {
		ctx   context.Context
		email string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *model.Customer
		wantErr bool
	}{
		{
			name: "best case",
			fields: fields{
				data: customer,
				err:  nil,
			},
			args: args{
				ctx:   context.Background(),
				email: "username@mail.com",
			},
			want:    &customer,
			wantErr: false,
		},
		{
			name: "GetContext Error",
			fields: fields{
				data: model.Customer{},
				err:  errors.New("some error"),
			},
			args: args{
				ctx:   context.Background(),
				email: "username@mail.com",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, m, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			sqlxDB := sqlx.NewDb(db, "sqlmock")
			testDB := postgresql.NewCustomerRepository(sqlxDB)

			query := "SELECT id, email, username, password FROM customers WHERE email = $1"

			mockExpectQuery := m.ExpectQuery(query).WithArgs(tt.args.email)
			if tt.fields.err != nil {
				mockExpectQuery.WillReturnError(err)
			} else {
				row := sqlmock.NewRows([]string{"id", "email", "username", "password"}).AddRow(tt.fields.data.ID, tt.fields.data.Email, tt.fields.data.Username, tt.fields.data.Password)

				mockExpectQuery.WillReturnRows(row)
			}

			got, err := testDB.GetCustomerByEmail(tt.args.ctx, tt.args.email)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.want, got)
		})
	}
}
