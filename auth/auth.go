package auth

import (
	"Cloud/dataBase"
	"Cloud/email"
	"Cloud/logger"
	"Cloud/models"
	"Cloud/utils"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"os"
	"time"
)

// Получаем секретный ключ из переменной окружения
var jwtKey = []byte(os.Getenv("JWT_SECRET_KEY"))

// Создание JWT токена
func GenerateAccessToken(user models.User) (string, time.Time, error) {
	expirationTime := time.Now().Add(15 * time.Minute) // Время действия токена - 15 минут
	claims := &models.Claims{
		Email:  user.Email,
		UserID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	return tokenString, expirationTime, err
}

// Проверка и валидация токена
func ValidateJWT(tokenStr string) (*models.Claims, error) {
	claims := &models.Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		logger.Error("Ошибка при парсинге токена: " + err.Error())
	}
	if !token.Valid {
		logger.Error("Токен недействителен")
		return nil, fmt.Errorf("токен недействителен")
	}

	return claims, nil
}

func GenerateRefreshToken(user models.User) (string, error) {
	expirationTime := time.Now().Add(30 * 24 * time.Hour) // Время действия рефреш токена — 30 дней
	claims := &models.Claims{
		Email: user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	return tokenString, err
}

func RefreshTokenHandler(w http.ResponseWriter, r *http.Request) {

	// Извлекаем рефреш токен из куки запроса
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		logger.Error("Отсутствует рефреш токен: " + err.Error())
		http.Error(w, "Отсутствует рефреш токен", http.StatusUnauthorized)
		return
	}

	// Валидация рефреш токена
	// Если токен недействителен или истёк, сервер снова возвращает ошибку 401 и логирует информацию о проблеме.
	claims, err := ValidateJWT(cookie.Value)
	if err != nil || claims == nil {
		logger.Error("Недействительный рефреш токен: " + err.Error())
		http.Error(w, "Недействительный рефреш токен", http.StatusUnauthorized)
		return
	}

	// Создание нового access токена
	// Если рефреш токен валиден, создаётся новый токен доступа с использованием email из рефреш токена. Функция GenerateJWT отвечает за создание нового
	// access токена. Если происходит ошибка во время генерации, сервер возвращает ошибку 500 (Internal Server Error) и логирует это событие.
	user := models.User{Email: claims.Email}
	accessToken, _, err := GenerateAccessToken(user)
	if err != nil {
		logger.Error("Ошибка генерации access токена: " + err.Error())
		http.Error(w, "Ошибка генерации токена", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"accessToken": accessToken})
}

func ConfirmEmailHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request struct {
			Email string `json:"email"`
			Code  string `json:"code"`
		}

		// Декодируем запрос
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			logger.Error("Ошибка декодирования JSON: " + err.Error())
			http.Error(w, "Неверный формат данных", http.StatusBadRequest)
			return
		}

		// Проверяем, существует ли код для данного email
		storedData, exists := models.TemporaryStore[request.Email]
		if !exists {
			logger.Error("Email не найден или код подтверждения просрочен: " + request.Email)
			http.Error(w, "Email не найден или код подтверждения просрочен", http.StatusNotFound)
			return
		}

		// Проверяем, истек ли срок действия кода подтверждения
		if time.Since(storedData.CreatedAt) > time.Hour {
			logger.Error("Код подтверждения для email просрочен: " + request.Email)
			delete(models.TemporaryStore, request.Email) // Удаляем просроченный код
			http.Error(w, "Код подтверждения просрочен", http.StatusUnauthorized)
			return
		}

		// Проверяем соответствие кода
		if storedData.Code != request.Code {
			logger.Error("Неверный код подтверждения для email: " + request.Email)
			http.Error(w, "Неверный код подтверждения", http.StatusUnauthorized)
			return
		}

		// Достаём данные пользователя
		user := storedData.User

		// Установка даты создания и обновления пользователя
		user.FromDateCreate = time.Now().Format(time.RFC3339)
		user.FromDateUpdate = user.FromDateCreate

		// Хеширование пароля перед сохранением
		user.Password, err = utils.HashPassword(user.Password)
		if err != nil {
			logger.Error("Ошибка хеширования пароля: " + err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Сохраняем пользователя в базе данных
		err = dataBase.DBCreateUser(db, &user)
		if err != nil {
			logger.Error("Ошибка создания пользователя: " + err.Error())
			http.Error(w, "Не удалось создать пользователя", http.StatusInternalServerError)
			return
		}

		// Удаляем код подтверждения из временного хранилища
		delete(models.TemporaryStore, request.Email)

		// Успешное подтверждение
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Email " + request.Email + " успешно подтвержден!"})
	}
}

func ResendConfirmationEmailHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request struct {
			Email string `json:"email"`
		}

		// Декодируем запрос
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			logger.Error("Ошибка декодирования JSON: " + err.Error())
			http.Error(w, "Неверный формат данных", http.StatusBadRequest)
			return
		}

		// Генерируем новый код подтверждения
		code := utils.GenRandCode()

		// Сохраняем новый код в TemporaryStore
		models.TemporaryStore[request.Email] = struct {
			Code      string
			User      models.User
			CreatedAt time.Time
		}{
			Code:      code,
			User:      models.User{Email: code},
			CreatedAt: time.Now(),
		}

		// Отправляем код на почту
		err = email.SendConfirmationEmail(request.Email, code)
		if err != nil {
			logger.Error("Ошибка отправки email: " + err.Error())
			http.Error(w, "Ошибка отправки email", http.StatusInternalServerError)
			return
		}

		// Успешная отправка
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode("Повторный код подтверждения отправлен на email " + request.Email)
	}
}
