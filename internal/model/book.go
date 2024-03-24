package model

import (
	"time"

	"github.com/lib/pq"
)

type Book struct {
	ID         uint        `db:"id"`
	Title      string      `db:"title"`
	Author     string      `db:"author"`
	Price      float64     `db:"price"`
	CategoryID uint        `db:"category_id"`
	CreatedAt  time.Time   `db:"created_at"`
	DeletedAt  pq.NullTime `db:"deleted_at"`
}
