package url

import (
	"context"
	"database/sql"

	"github.com/Neimess/shortener/internal/model/url"
)

type Repository struct{ db *sql.DB }

func New(db *sql.DB) url.Repository { return &Repository{db: db} }

func (r *Repository) Save(ctx context.Context, u *url.URL) error {
	const query = `
	INSERT INTO urls (
        id, original, short_code,  created_at, clicks, expires_at
    ) VALUES ($1, $2, $3, $4, $5, $6)
    `
	_, err := r.db.ExecContext(ctx, query,
		u.ID,
		u.Original,
		u.ShortCode,
		u.CreatedAt,
		u.Clicks,
		u.ExpiresAt)
	return err
}

func (r *Repository) FindByCode(ctx context.Context, code string) (*url.URL, error) {
	const query = `
        SELECT id, original, short_code, created_at, clicks, expires_at
        FROM urls
        WHERE short_code = $1
    `
	row := r.db.QueryRowContext(ctx, query, code)
	var u url.URL
	if err := row.Scan(
		&u.ID,
		&u.Original,
		&u.ShortCode,
		&u.CreatedAt,
		&u.Clicks,
		&u.ExpiresAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

func (r *Repository) DeleteExpired(ctx context.Context) (int64, error) {
	const query = `
                DELETE FROM urls
                WHERE expires_at IS NOT NULL
                 AND expires_at < NOW()`
	res, err := r.db.ExecContext(ctx, query)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}
