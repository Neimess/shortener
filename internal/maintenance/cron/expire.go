package cron

import (
	"context"
	"log"
	"time"

	repo "github.com/Neimess/shortener/internal/infrastructure/db/postgres/url"
)

func ExpiryCleaner(repo repo.Repository, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		cancel()

		rows, err := repo.DeleteExpired(ctx)
		if err != nil {
			log.Println("Error cleaning expired URLs:", err)
			continue
		}

		if rows > 0 {
			log.Printf("Cleaned %d expired URLs\n", rows)
		}
	}

}
