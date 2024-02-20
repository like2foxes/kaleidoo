package main

import (
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	m, err := migrate.New(
		"file://sql/migrations",
		"postgres://nirlir:nirlir@localhost:5432/nirlir?sslmode=disable",
	)

	if err != nil {
		log.Fatal(err)
	}
	if err = m.Up(); err != nil {
		log.Fatal(err)
	}
}
