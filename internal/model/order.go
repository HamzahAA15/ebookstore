package model

import (
	"time"

	"github.com/lib/pq"
)

type Order struct {
	ID                uint        `db:"id"`
	CustomerID        uint        `db:"customer_id"`
	CustomerReference string      `db:"customer_reference"`
	ReceiverName      string      `db:"receiver_name"`
	Address           string      `db:"address"`
	City              string      `db:"city"`
	District          string      `db:"district"`
	PostalCode        string      `db:"postal_code"`
	Shipper           string      `db:"shipper"`
	AirwaybillNumber  string      `db:"airwaybill_number"`
	OrderDate         time.Time   `db:"order_date"`
	TotalItem         int         `db:"total_item"`
	TotalPrice        float64     `db:"total_price"`
	UpdatedAt         pq.NullTime `db:"updated_at"`
	DeletedAt         pq.NullTime `db:"deleted_at"`
}

type Item struct {
	ID        uint        `db:"id"`
	BookID    uint        `db:"book_id"`
	Quantity  int         `db:"quantity"`
	OrderID   uint        `db:"order_id"`
	CreatedAt time.Time   `db:"created_at"`
	DeletedAt pq.NullTime `db:"deleted_at"`
}
