package main

import (
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	m, err := migrate.New(
		"file:///home/migrations/schema",
		os.Getenv("MIGRATIONS_CONNECTION_STRING"),
	)
	if err != nil {
		log.Fatal("can't run migrations: ", err)
		os.Exit(1)
	}

	m.Up()
}
