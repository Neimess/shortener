package main

import (
	"log"
	"net/http"
	"time"

	app "github.com/Neimess/shortener/internal/bootstrap"
	_ "github.com/Neimess/shortener/docs"
)

// @title        URL Shortener API
// @version      1.0
// @description  Сервис для сокращения и редиректа URL.
// @host         localhost:8080
// @BasePath     /
func main() {
	app := app.Initialize()

	srv := &http.Server{
		Addr:         ":" + app.Config.Port,
		Handler:      app.Handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("Server running on %s", srv.Addr)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server error: %v", err)
	}
}
