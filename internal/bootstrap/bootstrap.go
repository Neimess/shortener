package bootstrap

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/Neimess/shortener/internal/config"
	"github.com/Neimess/shortener/internal/infrastructure/cache"
	"github.com/Neimess/shortener/internal/infrastructure/db"
	"github.com/Neimess/shortener/internal/infrastructure/generator"
	authService "github.com/Neimess/shortener/internal/service/auth"
	urlService "github.com/Neimess/shortener/internal/service/url"
	jwtutil "github.com/Neimess/shortener/internal/util/jwt"

	authHandler "github.com/Neimess/shortener/internal/api/handler/auth"
	urlHandler "github.com/Neimess/shortener/internal/api/handler/url"
	"github.com/Neimess/shortener/internal/api/middleware"
	"github.com/Neimess/shortener/internal/api/router"
	"github.com/Neimess/shortener/internal/maintenance/metric"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type App struct {
	Config   config.Config
	ServeMux *http.ServeMux
}

func New() *App {
	cfg := config.New()
	metrics.Init()

	dbConn, err := sql.Open(cfg.Driver(), cfg.DSN())
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	dbConn.SetMaxOpenConns(30)
	dbConn.SetMaxIdleConns(10)
	dbConn.SetConnMaxLifetime(5 * time.Minute)

	repos := db.NewRepositories(dbConn, cfg)

	var cacheClient cache.FullCache

	cacheClient, err = cache.NewRedisAdapter(cfg.RedisAddr())
	if err != nil {
		log.Printf("redis init failed: %v — using NullCache", err)
		cacheClient = cache.NewNullAdapter()
	}

	// Вспомогательные функции
	codeGen := generator.NewCharsetGenerator()
	jwtManager := jwtutil.New(
		cfg.JWTSecret(),
		cfg.AccessTTL(),
		cfg.RefreshTTL(),
	)
	// Services
	urlSvc := urlService.NewService(repos, cacheClient, codeGen)
	authSvc := authService.NewService(repos, cacheClient, jwtManager)

	// Handlers
	urlH := urlHandler.NewURLHandler(urlSvc)
	authH := authHandler.NewAuthHandler(authSvc, jwtManager)

	// Routers
	mux := http.NewServeMux()
	mux.Handle("/metrics", middleware.LoggerMiddleware(middleware.PrometheusMiddleware(promhttp.Handler())))
	mux.Handle("/docs/", httpSwagger.Handler(
		httpSwagger.URL("/docs/doc.json"),
	))

	router := router.New(mux, cacheClient, jwtManager, 100, time.Minute)
	router.URLRoutes(urlH)
	router.AuthRoutes(authH)
	// router.UserRoutes(urlHandler)

	return &App{
		Config:   cfg,
		ServeMux: mux,
	}
}
