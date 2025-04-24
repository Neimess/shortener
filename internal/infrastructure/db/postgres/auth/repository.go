package auth

import (
	// "context"
	"context"
	"database/sql"

	"github.com/Neimess/shortener/internal/model/user"
)

type Repository struct {
	db *sql.DB
}

func New(db *sql.DB) user.Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, email, passwordHash string) (int, error) {
	const query = `
	INSERT INTO users (email, password_hash)
	VALUES ($1 $2)
	`
    var userID int
	if err := r.db.QueryRowContext(ctx, query, email, passwordHash).Scan(&userID); err != nil {
        return 0, nil
    }
    return userID, nil
}

func (r *Repository) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	const query = `
        SELECT id, email, password_hash, created_at, updated_at
        FROM users
        WHERE email = $1
    `
	u := &user.User{}
    err := r.db.QueryRowContext(ctx, query, email).Scan(
        &u.ID,
        &u.Email,
        &u.PasswordHash,
        &u.CreatedAt,
        &u.UpdatedAt,
    )
	if err == sql.ErrNoRows {
        return nil, nil
    }
	if err != nil {
        return nil, err
    }
    return u, nil
}
