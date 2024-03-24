package request

type CreateOrder struct {
	Items        []Item `json:"items"`
	ReceiverName string `json:"receiver_name"`
	Address      string `json:"address"`
	City         string `json:"city"`
	District     string `json:"district"`
	PostalCode   string `json:"postal_code"`
	Shipper      string `json:"shipper"`
}

type Item struct {
	BookID   uint    `json:"book_id"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
}
