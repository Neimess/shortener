package user

import (
    "github.com/Neimess/shortener/internal/validators"
)

type RegisterInput struct {
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=6"`
}

func (r *RegisterInput) Validate() error {
    return validator.Validate.Struct(r)
}
