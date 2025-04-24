package auth

// swagger:model RefreshRequest
type RefreshRequest struct {
    RefreshToken string `json:"refresh_token"`
}

// swagger:model RefreshResponse
type RefreshResponse struct {
    AccessToken  string `json:"access_token"`
    RefreshToken string `json:"refresh_token,omitempty"`
}


// swagger:model RegisterRequest
type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// swagger:model RegisterResponse
type RegisterResponse struct {
	UserID int `json:"user_id"`
}

// swagger:model LoginRequest
type LoginRequest struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}

// swagger:model LoginResponse
type LoginResponse struct {
    AccessToken  string `json:"access_token"`
    RefreshToken string `json:"refresh_token,omitempty"`
}