package postgresql_test

import (
	"context"
	"database/sql"
	"ebookstore/internal/model"
	"ebookstore/internal/repository/postgresql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func Test_orderRepository_GetItemsByOrderID(t *testing.T) {
	items := []model.Item{
		{
			ID:       1,
			BookID:   1,
			OrderID:  1,
			Quantity: 1,
		},
	}

	type fields struct {
		data []model.Item
		err  error
	}
	type args struct {
		ctx     context.Context
		orderID uint
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []model.Item
		wantErr bool
	}{
		{
			name: "best case",
			fields: fields{
				data: items,
				err:  nil,
			},
			args: args{
				ctx:     context.Background(),
				orderID: 1,
			},
			want:    items,
			wantErr: false,
		},
		{
			name: "QueryContext Error",
			fields: fields{
				data: nil,
				err:  errors.New("error"),
			},
			args: args{
				ctx:     context.Background(),
				orderID: 1,
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
			testDB := postgresql.NewOrderRepository(sqlxDB)

			queryString := `
			SELECT
				id,
				book_id,
				order_id,
				quantity
			FROM items
			WHERE order_id = $1`

			mockExpectQuery := m.ExpectQuery(queryString).WithArgs(tt.args.orderID)
			if tt.fields.err != nil {
				mockExpectQuery.WillReturnError(err)
			} else {
				row := sqlmock.NewRows([]string{"id", "book_id", "order_id", "quantity"}).AddRow(items[0].ID, items[0].BookID, items[0].OrderID, items[0].Quantity)

				mockExpectQuery.WillReturnRows(row)
			}

			got, err := testDB.GetItemsByOrderID(tt.args.ctx, tt.args.orderID)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_orderRepository_GetOrderHistoryByCustomerID(t *testing.T) {
	orders := []model.Order{
		{
			ID:                1,
			CustomerID:        1,
			CustomerReference: "ref",
			ReceiverName:      "name",
			Address:           "address",
			City:              "city",
			District:          "district",
			PostalCode:        "123",
			Shipper:           "shipper",
			AirwaybillNumber:  "123",
			OrderDate:         time.Now().UTC().Truncate(time.Minute),
			TotalItem:         1,
			TotalPrice:        1,
		},
	}

	type fields struct {
		data []model.Order
		err  error
	}
	type args struct {
		ctx        context.Context
		cusomterID uint
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []model.Order
		wantErr bool
	}{
		{
			name: "best case",
			fields: fields{
				data: orders,
				err:  nil,
			},
			args: args{
				ctx:        context.Background(),
				cusomterID: 1,
			},
			want:    orders,
			wantErr: false,
		},
		{
			name: "SelectContext error",
			fields: fields{
				data: nil,
				err:  errors.New("error"),
			},
			args: args{
				ctx:        context.Background(),
				cusomterID: 1,
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
			testDB := postgresql.NewOrderRepository(sqlxDB)

			query := `
				SELECT 
					id,
					customer_id,
					customer_reference,
					receiver_name,
					address,
					city,
					district,
					postal_code,
					shipper,
					airwaybill_number,
					order_date,
					total_item,
					total_price
				FROM orders 
				WHERE customer_id = $1 AND deleted_at is NULL
				ORDER By order_date DESC`

			mockExpectQuery := m.ExpectQuery(query).WithArgs(tt.args.cusomterID)
			if tt.fields.err != nil {
				mockExpectQuery.WillReturnError(err)
			} else {
				row := sqlmock.NewRows([]string{"id", "customer_id", "customer_reference", "receiver_name", "address", "city", "district", "postal_code", "shipper", "airwaybill_number", "order_date", "total_item", "total_price"})
				for _, order := range tt.fields.data {
					row.AddRow(
						order.ID,
						order.CustomerID,
						order.CustomerReference,
						order.ReceiverName,
						order.Address,
						order.City,
						order.District,
						order.PostalCode,
						order.Shipper,
						order.AirwaybillNumber,
						order.OrderDate,
						order.TotalItem,
						order.TotalPrice,
					)
				}

				mockExpectQuery.WillReturnRows(row)
			}

			got, err := testDB.GetOrderHistoryByCustomerID(tt.args.ctx, tt.args.cusomterID)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_orderRepository_CreateOrder(t *testing.T) {
	order := model.Order{
		ID:                1,
		CustomerID:        1,
		CustomerReference: "ref",
		ReceiverName:      "name",
		Address:           "address",
		City:              "city",
		District:          "district",
		PostalCode:        "123",
		Shipper:           "shipper",
		AirwaybillNumber:  "123",
		OrderDate:         time.Now().UTC().Truncate(time.Minute),
		TotalItem:         1,
		TotalPrice:        1,
	}

	type fields struct {
		data uint
		err  error
	}
	type args struct {
		ctx context.Context
		req model.Order
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
				data: order.ID,
				err:  nil,
			},
			args: args{
				ctx: context.Background(),
				req: order,
			},
			want:    order.ID,
			wantErr: false,
		},
		{
			name: "QueryRowxContext error",
			fields: fields{
				data: 0,
				err:  errors.New("error"),
			},
			args: args{
				ctx: context.Background(),
				req: order,
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
			testDB := postgresql.NewOrderRepository(sqlxDB)
			m.ExpectBegin()

			query := `
			INSERT INTO orders (
				customer_id,
				customer_reference,
				receiver_name,
				address,
				city,
				district,
				postal_code,
				order_date,
				shipper,
				airwaybill_number
			) VALUES (
				$1,
				$2,
				$3,
				$4,
				$5,
				$6,
				$7,
				$8,
				$9,
				$10
			) RETURNING id;`

			tx, _ := sqlxDB.BeginTxx(tt.args.ctx, &sql.TxOptions{Isolation: sql.LevelRepeatableRead})

			mockExpectQuery := m.ExpectQuery(query).WithArgs(tt.args.req.CustomerID, tt.args.req.CustomerReference, tt.args.req.ReceiverName, tt.args.req.Address, tt.args.req.City, tt.args.req.District, tt.args.req.PostalCode, tt.args.req.OrderDate, tt.args.req.Shipper, tt.args.req.AirwaybillNumber)
			if tt.fields.err != nil {
				mockExpectQuery.WillReturnError(err)
			} else {
				mockExpectQuery.WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(tt.args.req.ID))
			}

			got, err := testDB.CreateOrder(tt.args.ctx, tx, tt.args.req)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_orderRepository_CreateItem(t *testing.T) {
	item := model.Item{
		ID:        1,
		BookID:    1,
		OrderID:   1,
		Quantity:  1,
		CreatedAt: time.Now().UTC().Truncate(time.Minute),
	}
	type fields struct {
		err error
	}
	type args struct {
		ctx  context.Context
		item model.Item
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
				err: nil,
			},
			args: args{
				ctx:  context.Background(),
				item: item,
			},
			wantErr: false,
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
			testDB := postgresql.NewOrderRepository(sqlxDB)
			m.ExpectBegin()

			query := "INSERT INTO items (book_id,quantity,order_id,created_at) VALUES ($1,$2,$3,$4)"

			tx, _ := sqlxDB.BeginTxx(tt.args.ctx, &sql.TxOptions{Isolation: sql.LevelRepeatableRead})

			mockExpectExec := m.ExpectExec(query).WithArgs(tt.args.item.BookID, tt.args.item.Quantity, tt.args.item.OrderID, tt.args.item.CreatedAt)
			if tt.fields.err != nil {
				mockExpectExec.WillReturnError(tt.fields.err)
			} else {
				mockExpectExec.WillReturnResult(sqlmock.NewResult(1, 1))
			}

			err = testDB.CreateItem(tt.args.ctx, tx, tt.args.item)
			if err != nil {
				println("asd", err.Error())
			}
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}

func Test_orderRepository_UpdateOrderByOrderID(t *testing.T) {
	order := model.Order{
		ID:         1,
		TotalItem:  2,
		TotalPrice: 10,
	}

	type mockExec struct {
		err error
	}

	type fields struct {
		db *sqlx.DB
	}

	type args struct {
		ctx   context.Context
		order model.Order
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantErr  bool
		mockExec mockExec
	}{
		{
			name: "best case",
			fields: fields{
				db: func() *sqlx.DB {
					return &sqlx.DB{}
				}(),
			},
			args: args{
				ctx:   context.Background(),
				order: order,
			},
			mockExec: mockExec{
				err: nil,
			},
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
			o := postgresql.NewOrderRepository(sqlxDB)
			m.ExpectBegin()

			query := "UPDATE orders SET total_item = $1, total_price = $2 WHERE id = $3"

			mockExpectExect := m.ExpectExec(query).WithArgs(tt.args.order.TotalItem, order.TotalPrice, tt.args.order.ID)
			if tt.wantErr {
				mockExpectExect.WillReturnError(tt.mockExec.err)
			} else {
				mockExpectExect.WillReturnResult(sqlmock.NewResult(1, 1))
			}

			tx, _ := sqlxDB.BeginTxx(tt.args.ctx, &sql.TxOptions{Isolation: sql.LevelRepeatableRead})

			err = o.UpdateOrderByOrderID(context.Background(), tx, order)
			assert.Equal(t, tt.mockExec.err, err)
		})
	}
}
