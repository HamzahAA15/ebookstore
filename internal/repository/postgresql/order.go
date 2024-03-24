package repository

import (
	"context"
	"ebookstore/internal/model"
	"ebookstore/internal/repository"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type orderRepository struct {
	db *sqlx.DB
}

func NewOrderRepository(db *sqlx.DB) repository.IOrderRepository {
	return &orderRepository{db: db}
}

func (o *orderRepository) GetItemsByOrderID(ctx context.Context, customerID uint) ([]model.Item, error) {

	var items []model.Item
	queryString := `
		SELECT
			id,
			book_id,
			order_id,
			quantity
		FROM items
		WHERE order_id = $1`

	query, err := o.db.PrepareContext(ctx, queryString)
	if err != nil {
		return nil, err
	}
	defer query.Close()

	itemsRows, err := query.QueryContext(ctx, customerID)
	if err != nil {
		return nil, err
	}

	for itemsRows.Next() {
		var item model.Item
		if err := itemsRows.Scan(&item.ID, &item.BookID, &item.OrderID, &item.Quantity); err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil

}

func (o *orderRepository) CreateOrder(ctx context.Context, req model.Order) (uint, error) {
	var id uint
	queryStr := `
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
		:customer_id,
		:customer_reference,
		:receiver_name,
		:address,
		:city,
		:district,
		:postal_code,
		:order_date,
		:shipper,
		:airwaybill_number
	) RETURNING id;`

	query, err := o.db.PrepareNamedContext(ctx, queryStr)
	if err != nil {
		return 0, err
	}
	defer query.Close()

	err = query.QueryRowxContext(ctx, req).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to create order: %s", err.Error())
	}

	return id, nil
}

func (o *orderRepository) CreateItem(ctx context.Context, item model.Item) error {
	queryStr := `
	INSERT INTO items (
		book_id,
		quantity,
		order_id,
		created_at
	) VALUES (
		:book_id,
		:quantity,
		:order_id,
		:created_at
	)`

	query, err := o.db.PrepareNamedContext(ctx, queryStr)
	if err != nil {
		return err
	}
	defer query.Close()

	_, err = query.ExecContext(ctx, item)
	if err != nil {
		return err
	}

	return nil
}

func (o *orderRepository) GetOrderHistoryByCustomerID(ctx context.Context, cusomterID uint) ([]model.Order, error) {
	var orders []model.Order
	queryString := `
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

	query, err := o.db.PrepareContext(ctx, queryString)
	if err != nil {
		return nil, err
	}
	defer query.Close()

	ordersRows, err := query.QueryContext(ctx, cusomterID)
	if err != nil {
		return nil, err
	}

	for ordersRows.Next() {
		var order model.Order
		err = ordersRows.Scan(
			&order.ID,
			&order.CustomerID,
			&order.CustomerReference,
			&order.ReceiverName,
			&order.Address,
			&order.City,
			&order.District,
			&order.PostalCode,
			&order.Shipper,
			&order.AirwaybillNumber,
			&order.OrderDate,
			&order.TotalItem,
			&order.TotalPrice,
		)
		if err != nil {
			return nil, err
		}

		orders = append(orders, order)
	}

	return orders, nil
}

func (o *orderRepository) UpdateOrderByOrderID(ctx context.Context, order model.Order) error {
	query, err := o.db.PrepareContext(ctx, "UPDATE orders SET total_item = $1, total_price = $2 WHERE id = $3")
	if err != nil {
		return err
	}
	defer query.Close()

	_, err = query.ExecContext(ctx, order.TotalItem, order.TotalPrice, order.ID)
	if err != nil {
		return err
	}

	return nil
}
