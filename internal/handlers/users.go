package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"ozinshe/internal/middleware"
	"ozinshe/internal/repositories"
	"ozinshe/internal/services"
	"time"
)

type UserHandler struct {
	service *services.UserService
}

func NewUserHandler(service *services.UserService) *UserHandler {
	return &UserHandler{service: service}
}

type userRegister struct {
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	Password    string    `json:"password"`
	PhoneNumber string    `json:"phone_number"`
	BornAt      time.Time `json:"born_at"`
}

type userLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginResponse struct {
	Token string `json:"token"`
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req userRegister

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.Name == "" || req.Email == "" || req.Password == "" {
		http.Error(w, "Вы ввели пустое значение!", http.StatusBadRequest)
		return
	}

	err := h.service.Register(ctx, req.Name, req.Email, req.Password, req.PhoneNumber, req.BornAt)
	if err != nil {

		http.Error(w, "Ошибка при регистрации пользователя", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req userLogin

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.Password == "" {
		http.Error(w, "Вы ввели пустое значение!", http.StatusBadRequest)
		return
	}

	token, err := h.service.Login(ctx, req.Email, req.Password)
	if err != nil {
		if errors.Is(err, services.ErrInvalidCredentials) {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}

	response := loginResponse{
		Token: token,
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func (h *UserHandler) Me(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		http.Error(w, "Пользователь не найден", http.StatusUnauthorized)
		return
	}

	user, err := h.service.GetByID(r.Context(), userID)

	if err != nil {
		if errors.Is(err, repositories.ErrUserNotFound) {
			http.Error(w, "Пользователь не найден", http.StatusNotFound)
			return
		}

		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
