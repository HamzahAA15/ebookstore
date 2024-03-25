package main

import (
	"context"
	"ebookstore/internal/httpservice"
	"ebookstore/internal/repository"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "root"
	password = "password"
	dbname   = "book_db"
)

func main() {
	app := fiber.New()
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err := repository.ConnectPostgres(context.Background(), psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	httpservice.InitRoutes(app, db)

	app.Listen(":8080")
}
