package main

import (
	"errors"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	m, err := migrate.New(
		"file://migrations",
		databaseURL,
	)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		sourceErr, databaseErr := m.Close()

		if sourceErr != nil {
			log.Printf("failed to close migration source: %v", sourceErr)
		}
		if databaseErr != nil {
			log.Printf("failed to close database connection: %v", databaseErr)
		}
	}()

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Println("no migration changes")
			return
		}

		log.Fatal(err)
	}

	log.Println("migration completed")
}