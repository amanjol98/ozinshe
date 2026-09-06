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

type userUpdate struct {
	Name        string     `json:"name"`
	PhoneNumber string     `json:"phone_number"`
	BornAt      *time.Time `json:"born_at"`
}

type userAuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginResponse struct {
	Token string `json:"token"`
}

type updatePasswordRequest struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req userAuthRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.Password == "" {
		http.Error(w, "Вы ввели пустое значение!", http.StatusBadRequest)
		return
	}

	err := h.service.Register(ctx, req.Email, req.Password)
	if err != nil {
		if errors.Is(err, services.ErrEmailAlreadyExists) {
			http.Error(w, "email уже существует", http.StatusConflict)
			return
		}
		http.Error(w, "Ошибка при регистрации пользователя", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req userAuthRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.Password == "" {
		http.Error(w, "Вы ввели пустое значение!", http.StatusBadRequest)
		return
	}

	if len(req.Password) < 6 {
		http.Error(w, "Пароль должен содержать минимум 6 символов", http.StatusBadRequest)
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

// Me godoc
// @Summary Получить текщего пользователя
// @Description Возвращает данные пользователя, авторизованного через JWT.
// @Tags Users
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.User
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/me [get]
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

// Update godoc
// @Summary Обновить данные текущего пользователя
// @Description Обновляет имя, номер телефона и дату рождения текущего пользователя.
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body userUpdate true "Данные для обновления"
// @Success 200 {object} models.User
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/me [patch]
func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		http.Error(w, "Пользователь не найден", http.StatusUnauthorized)
		return
	}

	var req userUpdate

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Неверное тело запроса", http.StatusBadRequest)
		return
	}

	if req.Name == "" || req.PhoneNumber == "" {
		http.Error(w, "Вы ввели пустое значение!", http.StatusBadRequest)
		return
	}

	err := h.service.Update(r.Context(), userID, req.Name, req.PhoneNumber, req.BornAt)
	if err != nil {
		http.Error(w, "Не удалось обновить пользователя", http.StatusInternalServerError)
		return
	}

	user, err := h.service.GetByID(r.Context(), userID)
	if err != nil {
		http.Error(w, "Не удалось получить пользователя", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, "Не удалось отправить ответ", http.StatusInternalServerError)
		return
	}

}

// UpdatePassword godoc
// @Summary Обновить пароль текущего пользователя
// @Description Обновляет пароль текущего пользователя.
// @Tags Users
// @Accept json
// @Security BearerAuth
// @Param request body updatePasswordRequest true "Данные для обновления"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/me/password [patch]
func (h *UserHandler) UpdatePassword(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		http.Error(w, "Пользователь не найден", http.StatusUnauthorized)
		return
	}

	var req updatePasswordRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Неверное тело запроса", http.StatusBadRequest)
		return
	}

	if req.OldPassword == "" || req.NewPassword == "" {
		http.Error(w, "Старый и новый пароль обязательны", http.StatusBadRequest)
		return
	}

	if len(req.NewPassword) < 6 {
		http.Error(w, "Новый пароль должен содержать минимум 6 символов", http.StatusBadRequest)
		return
	}

	err := h.service.UpdatePassword(r.Context(), userID, req.OldPassword, req.NewPassword)
	if err != nil {
		if err.Error() == "старый пароль указан неверно" {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		http.Error(w, "Не удалось изменить пароль", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
