package jwtutil

import "github.com/golang-jwt/jwt/v5"

type Role string

const (
	RoleGuest Role = "guest"
	RoleUser  Role = "user"
	RoleAdmin Role = "admin"
)

type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   Role
	jwt.RegisteredClaims
}
