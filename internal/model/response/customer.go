package response

type Customer struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
	Token      string `json:"token,omitempty"`
}
