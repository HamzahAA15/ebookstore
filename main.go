package main

import (
	"context"
	"ebookstore/internal/httpservice"
	"ebookstore/internal/repository"
	"ebookstore/utils/config"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/tanimutomo/sqlfile"
)

func main() {
	app := fiber.New()
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable", config.Host, config.Port, config.User, config.Password, config.Dbname)

	db, err := repository.ConnectPostgres(context.Background(), psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	s := sqlfile.New()
	s.Directory("./db/migrations")
	s.Directory("./db/seeds")
	_, err = s.Exec(db.DB)
	if err != nil {
		log.Printf("Failed to apply migrations & seeder, run manually from sql scripts in ./db: %v", err)
	}

	httpservice.InitRoutes(app, db)

	app.Listen(":8080")
}
