package auth

import (
	"Cloud/dataBase"
	"Cloud/logger"
	"Cloud/models"
	"database/sql"
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

// Логика аутентификации пользователя
func LoginUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var loginReq LoginRequest
		var user *models.User
		var message string

		//Декодирование JSON из тела запроса
		err := json.NewDecoder(r.Body).Decode(&loginReq)
		if err != nil {
			logger.Error("Ошибка декодирования JSON: " + err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Проверяем, заполнены ли оба поля
		if loginReq.Email == "" && loginReq.Phone == "" {
			logger.Error("Не указаны email или телефон")
			http.Error(w, "Не указаны email или телефон", http.StatusBadRequest)
			return
		}

		// Проверяем какой из полей заполнен, и ищем пользователя
		if loginReq.Email != "" {
			user, message, err = dataBase.FindUserByEmail(db, loginReq.Email)
		}
		if loginReq.Phone != "" {
			user, message, err = dataBase.FindUserByPhone(db, loginReq.Phone)
		}

		// Проверяем на ошибки и статус пользователя
		if err != nil {
			logger.Error("Ошибка при поиске пользователя: " + err.Error())
			http.Error(w, "Ошибка при поиске пользователя", http.StatusInternalServerError)
			return
		}

		// Если пользователь заблокирован или удалён, выводим сообщение
		if message != "" {
			logger.Error(message)
			http.Error(w, message, http.StatusUnauthorized) // Ответ с соответствующим сообщением
		}

		// Если пользователь не найден стопаем функцию и выводим ответ
		if user == nil {
			return
		}

		// Проверка пароля
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginReq.Password)); err != nil {
			logger.Error("Неверный пароль")
			http.Error(w, "Неверный пароль", http.StatusUnauthorized)
			return
		}

		// Генерация токенов
		accessToken, err := GenerateAccessToken(*user)
		if err != nil {
			logger.Error("Ошибка создания токена: " + err.Error())
			http.Error(w, "Ошибка создания токена", http.StatusInternalServerError)
			return
		}

		refreshToken, err := GenerateRefreshToken(*user)
		if err != nil {
			logger.Error("Ошибка генерации refresh токена: " + err.Error())
			http.Error(w, "Ошибка генерации refresh токена", http.StatusInternalServerError)
			return
		}

		// Сохранение refresh токена в куки
		http.SetCookie(w, &http.Cookie{
			Name:     "refresh_token",
			Value:    refreshToken,
			Expires:  time.Now().Add(30 * 24 * time.Hour), // Время жизни куки
			HttpOnly: false,                               // Защита от доступа через JavaScript
			Secure:   false,                               // Убедитесь, что вы используете HTTPS
			Path:     "/",                                 // Путь для куки
		})

		// Возвращаем access токен
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"access_token": accessToken})
	}
}
