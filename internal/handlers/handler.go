package handlers

import (
	"AuthService/internal/domain"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

type Service interface {
	Create(context.Context, domain.User) error
	Auth(context.Context, domain.User) (string, error)
	Validate(string) error
}

type Handler struct {
	service Service
}

func NewHandler(s Service) *Handler {
	return &Handler{ service: s }
}

func (h *Handler) Auth(w http.ResponseWriter, r *http.Request) {
	var user domain.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		domain.FormError(err.Error(), time.Now(), w, 400)
		return
	}
	if err := user.ValidateUser(); err != nil {
		domain.FormError(err.Error(), time.Now(), w, 400)
		return
	}
	tokenString, err := h.service.Auth(r.Context(), user)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			domain.FormError(err.Error(), time.Now(), w, 404)
		} else {
			domain.FormError(err.Error(), time.Now(), w, 500)
		}
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name: "accessToken",
		Value: tokenString,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Path: "/",
		MaxAge: 15 * 60 * 60,
	})
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var user domain.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		domain.FormError(err.Error(), time.Now(), w, 400)
		return
	}
	if err := user.ValidateUser(); err != nil {
		domain.FormError(err.Error(), time.Now(), w, 400)
		return
	}
	if err := h.service.Create(r.Context(), user); err != nil {
		if errors.Is(err, domain.ErrUserAlreadyExists) {
			domain.FormError(err.Error(), time.Now(), w, 409)
		} else {
			domain.FormError(err.Error(), time.Now(), w, 500)
		}
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) Validate(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("accessToken")
	if err != nil {
		domain.FormError(err.Error(), time.Now(), w, 401)
		return
	}
	tokenString := cookie.Value
	if tokenString == "" {
		domain.FormError(errors.New("No token").Error(), time.Now(), w, 401)
		return
	}
	if err := h.service.Validate(tokenString); err != nil {
		domain.FormError(err.Error(), time.Now(), w, 401)
		return
	}
	w.WriteHeader(200)
}