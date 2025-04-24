package auth

import (
	// "context"
	"database/sql"

	"github.com/Neimess/shortener/internal/model/auth"
)

type Repository struct{ db *sql.DB }

func New(db *sql.DB) auth.Repository { return &Repository{db: db} }
