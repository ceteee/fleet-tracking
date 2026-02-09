package database

import (
	"database/sql"
	"log"
	"time"

	"fleet-management-system/internal/config"

	_ "github.com/lib/pq"
)

func ConnectPostgres(cfg *config.Config) *sql.DB {
	db, err := sql.Open("postgres", cfg.PostgresDSN())
	if err != nil {
		log.Fatal("failed to open database:", err)
	}

	const maxRetry = 10
	for i := 1; i <= maxRetry; i++ {
		err = db.Ping()
		if err == nil {
			log.Println("PostgreSQL connected")
			return db
		}

		log.Printf("PostgreSQL not ready, retry %d/%d...", i, maxRetry)
		time.Sleep(2 * time.Second)
	}

	log.Fatal("failed to connect database after retries:", err)
	return nil
}
