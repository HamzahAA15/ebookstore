package response

import "time"

type Order struct {
	StatusCode int             `json:"status_code"`
	Message    string          `json:"message"`
	Data       CreateOrderData `json:"data,omitempty"`
}

type CreateOrderData struct {
	OrderID           uint   `json:"order_id,omitempty"`
	CustomerReference string `json:"customer_reference,omitempty"`
	AirwaybillNumber  string `json:"airwaybill_number,omitempty"`
	OrderDate         string `json:"order_date,omitempty"`
}

type GetUserOrders struct {
	StatusCode int         `json:"status_code"`
	Message    string      `json:"message"`
	Data       []OrderData `json:"data,omitempty"`
}

type OrderData struct {
	OrderID           uint      `json:"order_id"`
	CustomerReference string    `json:"customer_reference"`
	ReceiverName      string    `json:"receiver_name"`
	Address           string    `json:"address"`
	City              string    `json:"city"`
	District          string    `json:"district"`
	PostalCode        string    `json:"postal_code"`
	Shipper           string    `json:"shipper"`
	AirwaybillNumber  string    `json:"airwaybill_number"`
	OrderDate         time.Time `json:"order_date"`
	Items             []Item    `json:"items"`
	TotalItem         int       `json:"total_item"`
	TotalPrice        float64   `json:"total_price"`
}

type Item struct {
	BookID   uint    `json:"book_id"`
	Title    string  `json:"title"`
	Author   string  `json:"author"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
}
