package model

type Customer struct {
	ID       uint   `db:"id"`
	Email    string `db:"email"`
	Password string `db:"password"`
	Username string `db:"username"`
}

const CustomerBodyEmailTemplate = `
<p>Dear %s,</p>
<p>Welcome to our online bookstore!</p>
<p>Your account has been successfully created.</p>
<p>Start exploring our collection of books and enjoy shopping with us!</p>
<p>Best regards,</p>
<p>The Bookstore Team</p>
`
