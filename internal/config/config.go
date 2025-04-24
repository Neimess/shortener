package config

import (
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"
)

type Config interface {
	ServerPort() string
	Driver() string
	DSN() string
	RedisAddr() string
	RedisPassword() string
	RedisDB() int
	JWTSecret() []byte
	AccessTTL() time.Duration
	RefreshTTL() time.Duration
}

type configImpl struct {
	port       string
	driver     string
	dsn        string
	redisAddr  string
	redisPass  string
	redisDB    int
	jwtSecret  []byte
	accessTTL  time.Duration
	refreshTTL time.Duration
}

var (
	instance Config
	once     sync.Once
)

// Load загружает конфигурацию приложения один раз.
//
// Использует переменные окружения для настройки параметров подключения к базе данных,
// Redis и веб-серверу. Если переменные окружения отсутствуют — применяются значения по умолчанию.
//
// Повторный вызов функции всегда возвращает уже загруженную конфигурацию (реализовано через sync.Once).
//
// Пример переменных окружения:
//
//	DRIVER=postgres
//	REDIS_ADDR=localhost:6379
func New() Config {
	once.Do(func() {
		// Основные настройки
		driver := getEnv("DRIVER", "postgres")
		dsn := buildDSN(driver)
		port := getEnv("PORT", "8080")

		// Redis
		redisAddr := getEnv("REDIS_ADDR", "localhost:6379")
		redisPass := getEnv("REDIS_PASSWORD", "")
		redisDB := getEnvInt("REDIS_DB", 0)

		// JWT
		secret := getEnv("JWT_SECRET", "secret")

		// TTL
		accessTTL := getEnvDuration("ACCESS_TTL", 15*time.Minute)
		refreshTTL := getEnvDuration("REFRESH_TTL", 7*24*time.Hour)

		instance = &configImpl{
			port:       port,
			driver:     driver,
			dsn:        dsn,
			redisAddr:  redisAddr,
			redisPass:  redisPass,
			redisDB:    redisDB,
			jwtSecret:  []byte(secret),
			accessTTL:  accessTTL,
			refreshTTL: refreshTTL,
		}
	})
	return instance
}

func (c *configImpl) ServerPort() string        { return c.port }
func (c *configImpl) Driver() string            { return c.driver }
func (c *configImpl) DSN() string               { return c.dsn }
func (c *configImpl) RedisAddr() string         { return c.redisAddr }
func (c *configImpl) RedisPassword() string     { return c.redisPass }
func (c *configImpl) RedisDB() int              { return c.redisDB }
func (c *configImpl) JWTSecret() []byte         { return c.jwtSecret }
func (c *configImpl) AccessTTL() time.Duration  { return c.accessTTL }
func (c *configImpl) RefreshTTL() time.Duration { return c.refreshTTL }

func buildDSN(driver string) string {
	switch driver {
	case "postgres":
		return buildPostgresDSN()
	case "sqlite3":
		return buildSQLiteDSN()
	default:
		panic("unsupported DB_DRIVER: " + driver)
	}
}

func buildPostgresDSN() string {
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "postgres")
	pass := getEnv("DB_PASSWORD", "")
	name := getEnv("DB_NAME", "shortener")
	ssl := getEnv("DB_SSLMODE", "disable")
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		user, pass, host, port, name, ssl,
	)
}

func buildSQLiteDSN() string {
	path := getEnv("DB_PATH", "shortener.db")
	return fmt.Sprintf("file:%s?_foreign_keys=1", path)
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return fallback
}

func getEnvDuration(key string, fallback time.Duration) time.Duration {
	if v := os.Getenv(key); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}
	return fallback
}
