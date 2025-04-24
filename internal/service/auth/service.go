package auth

import "context"

type Service interface {
	Register(ctx context.Context, email, password string) (int, error)
	Login(ctx context.Context, email, password string) (accessToken, refreshToken string, err error)
	Refresh(ctx context.Context, refreshToken string) (newAccessToken string, err error)
}
