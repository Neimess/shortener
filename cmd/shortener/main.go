package main

import (
	"log"
	"net/http"
	"time"

	_ "github.com/Neimess/shortener/docs"
	app "github.com/Neimess/shortener/internal/bootstrap"
)

// @title        URL Shortener API
// @version      1.0
// @description  Сервис для сокращения и редиректа URL.
// @host         localhost:8080
// @BasePath     /
func main() {
	app := app.New()

	srv := &http.Server{
		Addr:         ":" + app.Config.ServerPort(),
		Handler:      app.ServeMux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("Server running on %s", srv.Addr)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server error: %v", err)
	}
}
