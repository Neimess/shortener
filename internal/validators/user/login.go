package user

import (
	"github.com/Neimess/shortener/internal/validators"
)

type LoginInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func (l *LoginInput) Validate() error {
	return validator.Validate.Struct(l)
}
