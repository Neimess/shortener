package bootstrap

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"

	"github.com/Neimess/shortener/internal/config"
	"github.com/Neimess/shortener/internal/infrastructure/cache"
	"github.com/Neimess/shortener/internal/infrastructure/db"
	"github.com/Neimess/shortener/internal/infrastructure/generator"
	urlUsecase "github.com/Neimess/shortener/internal/service/url"

	// authUsecase "github.com/Neimess/shortener/internal/service/auth"
	urlHandler "github.com/Neimess/shortener/internal/api/handler/url"
	"github.com/Neimess/shortener/internal/api/middleware"
	"github.com/Neimess/shortener/internal/api/router"
	"github.com/prometheus/client_golang/prometheus/promhttp"

)

type App struct {
	Config  *config.Config
	Handler http.Handler
}


func Initialize() *App {
	cfg := config.Load()
	dbConn, err := sql.Open(cfg.Driver, cfg.DSN)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	dbConn.SetMaxOpenConns(30)
	dbConn.SetMaxIdleConns(10)
	dbConn.SetConnMaxLifetime(5 * time.Minute)

	repos := db.NewRepositories(dbConn, cfg)

	var cacheClient cache.FullCache
	cacheClient, err = cache.NewRedisAdapter(cfg.RedisAddr, cfg.RedisPassword, cfg.RedisDB)
	if err != nil {
		log.Printf("redis init failed: %v — using NullCache", err)
		cacheClient = cache.NewNullAdapter()
	}

	// Вспомогательные функции
	codeGen := generator.NewCharsetGenerator()

	// Services
	urlSvc := urlUsecase.NewService(repos, cacheClient, codeGen)
	// authSvc := authUsecase.NewService(repos)

	// Handlers
	urlH := urlHandler.NewURLHandler(urlSvc)
	// authH := handler.NewAuthHandler(authSvc)

	// Routers
	mux := http.NewServeMux()
	router.RegisterURLRoutes(mux, urlH, cacheClient, 10, time.Minute)
	// router.RegisterAuthRoutes(mux, authH, cacheClient, 10, time.Minute)
	mux.Handle("/metrics", middleware.LoggerMiddleware(middleware.PrometheusMiddleware(promhttp.Handler())))
	

	return &App{
		Config:  cfg,
		Handler: mux,
	}
}
