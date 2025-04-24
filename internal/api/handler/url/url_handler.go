package url

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/Neimess/shortener/internal/util/http"
	model "github.com/Neimess/shortener/internal/model/url"
	service "github.com/Neimess/shortener/internal/service/url"
)

type URLHandler interface {
	Shorten(w http.ResponseWriter, r *http.Request)
	Redirect(w http.ResponseWriter, r *http.Request)
	Health(w http.ResponseWriter, r *http.Request)
}

type urlHandler struct {
	svc service.Service
}

func NewURLHandler(s service.Service) URLHandler {
	return &urlHandler{svc: s}
}


// @Summary     Сокращение URL
// @Description Принимает JSON {"url": "..."} и возвращает короткий код
// @Tags        urls
// @Accept      json
// @Produce     json
// @Param       input body model.ShortenRequest true "Original URL"
// @Success     200 {object} model.ShortenResponse
// @Failure     400 {string} string "invalid payload"
// @Failure     500 {string} string "failed to shorten url"
// @Router      /shorten [post]
func (h *urlHandler) Shorten(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		httputil.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	var req model.ShortenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid payload")
		return
	}
	code, err := h.svc.Shorten(req.URL)
	if err != nil {
		log.Printf("failed to shorten URL %s: %v", req.URL, err)
		httputil.Error(w, http.StatusInternalServerError, "failed to shorten url")
		return
	}
	httputil.JSON(w, http.StatusCreated, model.ShortenResponse{ShortCode: code})

}

// Redirect обрабатывает переход по короткой ссылке
// @Summary     Редирект по короткому коду
// @Description Находит оригинальный URL по короткому коду и делает 302 редирект
// @Tags        urls
// @Produce     plain
// @Param       code path string true "Short URL code"
// @Success     302 {string} string "Redirect to original URL"
// @Failure     404 {string} string "Short code not found"
// @Router      /{code} [get]
func (h *urlHandler) Redirect(w http.ResponseWriter, r *http.Request) {
	code := strings.TrimPrefix(r.URL.Path, "/")
	if code == "" || strings.Contains(code, "/") {
		http.NotFound(w, r)
		return
	}
	original, err := h.svc.Resolve(code)
	if err != nil || original == "" {
		http.NotFound(w, r)
		return
	}
	http.Redirect(w, r, original, http.StatusFound)
}

// Health нужен для healthcheck
// @Summary     Health используется для начального состояния проекта
// @Description Возвращает 200 OK, если сервис работает
// @Tags        health
// @Produce     plain
// @Success     200 {string} string "ok"
// @Router      /healthz [get]
func (h *urlHandler) Health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("ok"))
	if err != nil {
		http.Error(w, "service is not ok", http.StatusInternalServerError)
		return
	}
}
