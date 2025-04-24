package url

import (
	"net/http"

	authModel "github.com/Neimess/shortener/internal/model/auth"
	authService "github.com/Neimess/shortener/internal/service/auth"
	httputil "github.com/Neimess/shortener/internal/util/http"
	"github.com/Neimess/shortener/internal/util/jwt"
)

type AuthHandler interface {
	Register(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
	Refresh(w http.ResponseWriter, r *http.Request)
}

type authHandler struct {
	svc authService.Service
}

func NewAuthHandler(svc authService.Service, jwt jwtutil.JWTManager) AuthHandler {
	return &authHandler{svc: svc}
}

// @Summary     Регистрация пользователя
// @Description Принимает JSON {"email":"...","password":"..."} и создаёт нового пользователя
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       input body authModel.RegisterRequest true "Credentials"
// @Success     201 {object} authModel.RegisterResponse
// @Failure     400 {string} string "invalid payload"
// @Failure     500 {string} string "failed to shorten url"
// @Router      /register [post]
func (h *authHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		httputil.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	var req authModel.RegisterRequest
	if err := httputil.Bind(r, &req); err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid payload")
		return
	}
	id, err := h.svc.Register(r.Context(), req.Email, req.Password)
	if err != nil {
		httputil.Error(w, http.StatusInternalServerError, "registration failed")
		return
	}
	httputil.JSON(w, http.StatusCreated, authModel.RegisterResponse{UserID: id})
}

// @Summary     Авторизация и получение токенов
// @Description Принимает json, проверяет базу, если все успешно возвращает JWT
// @Tags        urls
// @Produce     plain
// @Param       code path string true "Short URL code"
// @Success     302 {string} string "Redirect to original URL"
// @Failure     404 {string} string "Short code not found"
// @Router      /{code} [get]
func (h *authHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		httputil.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req authModel.LoginRequest
	if err := httputil.Bind(r, &req); err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid payload")
		return
	}

	accessToken, refreshToken, err := h.svc.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		httputil.Error(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	// http.SetCookie(w, &http.Cookie{
	// 	Name:     "refresh_token",
	// 	Value:    refreshToken,
	// 	Path:     "/",
	// 	Expires:  time.Now().Add(h.jwt.RefreshTTL()),
	// 	HttpOnly: true,
	// 	Secure:   true,
	// 	SameSite: http.SameSiteStrictMode,
	// })

	resp := authModel.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	httputil.JSON(w, http.StatusOK, resp)

}

// @Summary     Обновление access token
// @Description Принимает JSON {"refresh_token":"..."} и возвращает новые access и refresh токены
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       input body authModel.RefreshRequest true "Refresh Token"
// @Success     200 {object} authModel.RefreshResponse
// @Failure     400 {string} string "invalid payload"
// @Failure     401 {string} string "invalid refresh token"
// @Failure     500 {string} string "refresh failed"
// @Router      /refresh [post]
func (h *authHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		httputil.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req authModel.RefreshRequest
	if err := httputil.Bind(r, &req); err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid payload")
		return
	}

	newAT, err := h.svc.Refresh(r.Context(), req.RefreshToken)
	if err != nil {
		httputil.Error(w, http.StatusUnauthorized, "invalid refresh token")
		return
	}

	resp := authModel.RefreshResponse{
		AccessToken: newAT,
	}
	httputil.JSON(w, http.StatusOK, resp)
}
