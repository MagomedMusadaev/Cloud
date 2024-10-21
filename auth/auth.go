package auth

import (
	"Cloud/logger"
	"Cloud/models"
	"encoding/json"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"os"
	"time"
)

// Получаем секретный ключ из переменной окружения
var jwtKey = []byte(os.Getenv("JWT_SECRET_KEY"))

// Создание JWT токена
func GenerateAccessToken(user models.User) (string, error) {
	expirationTime := time.Now().Add(15 * time.Minute) // Время действия токена - 15 минут
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

// Проверка и валидация токена
func ValidateJWT(tokenStr string) (*models.Claims, error) {
	claims := &models.Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil || token.Valid {
		return nil, err
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
	cookie, err := r.Cookie("token")
	if err != nil {
		logger.Error("Отсутствует рефреш токен: " + err.Error())
		http.Error(w, "Отсутствует рефреш токен", http.StatusUnauthorized)
		return
	}

	// Валидация рефреш токена
	// Если токен недействителен или истёк, сервер снова возвращает ошибку 401 и логирует информацию о проблеме.
	claims, err := ValidateJWT(cookie.Value)
	if err != nil {
		logger.Error("Недействительный рефреш токен: " + err.Error())
		http.Error(w, "Недействительный рефреш токен", http.StatusUnauthorized)
		return
	}

	// Создание нового access токена
	// Если рефреш токен валиден, создаётся новый токен доступа с использованием email из рефреш токена. Функция GenerateJWT отвечает за создание нового
	// access токена. Если происходит ошибка во время генерации, сервер возвращает ошибку 500 (Internal Server Error) и логирует это событие.
	user := models.User{Email: claims.Email}
	accessToken, err := GenerateAccessToken(user)
	if err != nil {
		logger.Error("Ошибка генерации access токена: " + err.Error())
		http.Error(w, "Ошибка генерации токена", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"accessToken": accessToken})
}