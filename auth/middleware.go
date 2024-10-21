package auth

import (
	"Cloud/logger"
	"net/http"
	"strings"
)

// Мидлвар для проверки токена
func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Получение заголовка авторизации
		authHeader := r.Header.Get("Authorization")

		// Проверка наличия токена
		if authHeader == "" {
			logger.Error("Отсутствует токен авторизации")
			http.Error(w, "Требуется авторизация", http.StatusUnauthorized)
			return
		}

		// Извлечение токена
		tokenStr := strings.Split(authHeader, "Bearer ")[1]

		// Валидация токена
		claims, err := ValidateJWT(tokenStr)
		if err != nil {
			logger.Error("Недействительный токен: " + err.Error())
			http.Error(w, "Недействительный токен", http.StatusUnauthorized)
			return
		}

		// Добавление данных пользователя в контекст
		r.Header.Set("userEmail", claims.Email)

		// Передача управления следующему обработчику
		next.ServeHTTP(w, r)
	})
}
