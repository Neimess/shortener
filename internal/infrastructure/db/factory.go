package db

import (
	"database/sql"

	"github.com/Neimess/shortener/internal/config"

	pgAuth "github.com/Neimess/shortener/internal/infrastructure/db/postgres/auth"
	pgURL "github.com/Neimess/shortener/internal/infrastructure/db/postgres/url"
	urlModel "github.com/Neimess/shortener/internal/model/url"
	userModel "github.com/Neimess/shortener/internal/model/user"
	// sqAuth "github.com/Neimess/shortener/internal/infrastructure/db/sqlite/auth"
	// sqURL "github.com/Neimess/shortener/internal/infrastructure/db/sqlite/url"
)

type RepositorySet struct {
	URL  urlModel.Repository
	User userModel.Repository
}

func NewRepositories(db *sql.DB, cfg config.Config) RepositorySet {
	switch cfg.Driver() {
	case "postgres":
		return RepositorySet{
			URL:  pgURL.New(db),
			User: pgAuth.New(db),
		}
	default:
		panic("unsupported DB driver: " + cfg.Driver())
	}
}
