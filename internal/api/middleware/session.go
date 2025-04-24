package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/Neimess/shortener/internal/util/jwt"
)

type ctxKey string

const UserClaimsKey ctxKey = "userClaims"

func AuthMiddleware(jwtMgr jwtutil.JWTManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			auth := r.Header.Get("Authorization")
			if !strings.HasPrefix(auth, "Bearer ") {
				http.Error(w, "missing or invalid auth header", http.StatusUnauthorized)
				return
			}
			tokenStr := strings.TrimPrefix(auth, "Bearer ")

			claims, err := jwtMgr.Decode(tokenStr)
			if err != nil {
				http.Error(w, "invalid or expired token", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), UserClaimsKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func FromContext(ctx context.Context) (*jwtutil.Claims, bool) {
	c, ok := ctx.Value(UserClaimsKey).(*jwtutil.Claims)
	return c, ok
}
