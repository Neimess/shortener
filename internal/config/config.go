package config

import (
	"fmt"
	"os"
	"strconv"
	"sync"
)

type Config struct {
	Driver        string
	DSN           string
	RedisAddr     string
	RedisPassword string
	RedisDB       int
	Port          string
}

var (
	cfg  *Config
	once sync.Once
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
//	internal_PORT=8080
//
// Возвращает указатель на структуру конфигурации *Config.
func Load() *Config {
	once.Do(func() {
		driver := getEnv("DRIVER", "postgres")
		dsn := buildDSN(driver)

		cfg = &Config{
			Driver:        driver,
			DSN:           dsn,
			RedisAddr:     getEnv("REDIS_ADDR", "localhost:6379"),
			RedisPassword: getEnv("REDIS_PASSWORD", ""),
			RedisDB:       getEnvInt("REDIS_DB", 0),
			Port:          getEnv("internal_PORT", "8080"),
		}
	})
	return cfg
}

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
