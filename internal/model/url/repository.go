package url

import "context"

type Repository interface {
	Save(ctx context.Context, u *URL) error
	FindByCode(ctx context.Context, code string) (*URL, error)
	DeleteExpired(ctx context.Context) (int64, error)
}
