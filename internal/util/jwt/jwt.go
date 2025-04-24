package jwtutil

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type JWTManager interface {
	GenerateAccessToken(userID, email string) (string, error)
	GenerateRefreshToken(userID, email string) (string, error)
	Decode(tokenStr string) (*Claims, error)
	RefreshTTL() time.Duration
}
type jwtManager struct {
	secret     []byte
	parser     *jwt.Parser
	accessTTL  time.Duration
	refreshTTL time.Duration
}

func New(secret []byte, accessTTL, refreshTTL time.Duration) JWTManager {
	return &jwtManager{
		secret: secret,
		parser: jwt.NewParser(
			jwt.WithValidMethods([]string{"HS256"}),
			jwt.WithStrictDecoding(),
		),
		accessTTL:  accessTTL,
		refreshTTL: refreshTTL,
	}
}

func (j *jwtManager) GenerateAccessToken(userID, email string) (string, error) {
	return j.encode(userID, email, j.accessTTL)
}

func (j *jwtManager) GenerateRefreshToken(userID, email string) (string, error) {
	return j.encode(userID, email, j.refreshTTL)
}

func (j *jwtManager) Decode(tokenStr string) (*Claims, error) {
	claims := &Claims{}
	token, err := j.parser.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return j.secret, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, jwt.ErrTokenExpired
	}

	return claims, nil
}

func (j *jwtManager) RefreshTTL() time.Duration {
	return j.refreshTTL
}
func (j *jwtManager) encode(userID, email string, expiry time.Duration) (string, error) {
	claims := &Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secret)
}
