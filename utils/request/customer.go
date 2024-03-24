package request

type Register struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Username string `json:"username"`
}

type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
