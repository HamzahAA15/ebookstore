package postgresql

import (
	"context"
	"ebookstore/internal/model"
	"ebookstore/internal/repository"
	"ebookstore/utils/transactioner"

	"github.com/jmoiron/sqlx"
)

type orderRepository struct {
	db *sqlx.DB
}

func NewOrderRepository(db *sqlx.DB) repository.IOrderRepository {
	return &orderRepository{db: db}
}

func (o *orderRepository) GetItemsByOrderID(ctx context.Context, OrderID uint) ([]model.Item, error) {

	var items []model.Item
	queryString := ` 
		SELECT
			id,
			book_id,
			order_id,
			quantity
		FROM items
		WHERE order_id = $1`

	itemsRows, err := o.db.QueryContext(ctx, queryString, OrderID)
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

func (o *orderRepository) CreateItem(ctx context.Context, tx transactioner.TxxProvider, item model.Item) error {
	query := "INSERT INTO items (book_id,quantity,order_id,created_at) VALUES ($1,$2,$3,$4)"

	_, err := tx.ExecContext(ctx, query, item.BookID, item.Quantity, item.OrderID, item.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (o *orderRepository) CreateOrder(ctx context.Context, tx transactioner.TxxProvider, req model.Order) (uint, error) {
	var id uint
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

	var args = []interface{}{
		req.CustomerID,
		req.CustomerReference,
		req.ReceiverName,
		req.Address,
		req.City,
		req.District,
		req.PostalCode,
		req.OrderDate,
		req.Shipper,
		req.AirwaybillNumber,
	}

	err := tx.QueryRowContext(ctx, query, args...).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (o *orderRepository) GetOrderHistoryByCustomerID(ctx context.Context, cusomterID uint) ([]model.Order, error) {
	var orders []model.Order
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

	err := o.db.SelectContext(ctx, &orders, query, cusomterID)
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (o *orderRepository) UpdateOrderByOrderID(ctx context.Context, tx transactioner.TxxProvider, order model.Order) error {
	query := "UPDATE orders SET total_item = $1, total_price = $2 WHERE id = $3"
	_, err := tx.ExecContext(ctx, query, order.TotalItem, order.TotalPrice, order.ID)
	if err != nil {
		return err
	}

	return nil
}
