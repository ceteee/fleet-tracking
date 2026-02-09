package main

import (
	"fleet-management-system/internal/config"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
)

func main() {
	if os.Getenv("APP_ENV") == "" {
		_ = godotenv.Load()
	}

	cfg := config.Load()

	if len(os.Args) < 2 {
		log.Fatal("usage: migrate [up|dowm]")
	}

	migrator, err := migrate.New(
		"file://migrations",
		cfg.PostgresDSN(),
	)

	if err != nil {
		log.Fatal(err)
	}

	switch os.Args[1] {
	case "up":
		if err := migrator.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		}
	case "down":
		if err := migrator.Down(); err != nil {
			log.Fatal(err)
		}
	default:
		log.Fatal("Uknown command")
	}

	log.Println("migration finished")
}
