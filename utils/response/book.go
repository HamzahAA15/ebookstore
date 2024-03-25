package response

type GetBooks struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
	Data       []Book `json:"data"`
}

type Book struct {
	ID       uint    `json:"id"`
	Title    string  `json:"title"`
	Author   string  `json:"author"`
	Price    float64 `json:"price"`
	Category string  `json:"category"`
}
