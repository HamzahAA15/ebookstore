package model

type Customer struct {
	ID       uint   `db:"id"`
	Email    string `db:"email"`
	Password string `db:"password"`
	Username string `db:"username"`
}
