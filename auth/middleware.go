package auth

import (
	"Cloud/dataBase"
	"Cloud/internal"
	"Cloud/logger"
	"Cloud/utils"
	"database/sql"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// ResponseWriterWrapper добавляет возможность перехватывать статус ответа
type ResponseWriterWrapper struct {
	http.ResponseWriter
	StatusCode int
}

// Мидлвар для проверки токена
func JWTMiddleware(db *sql.DB, next http.Handler) http.Handler {
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

		// Получение времени истечения токена из базы данных
		tokenExpiration, err := dataBase.GetTokenExpiration(db, claims.UserID)
		if err != nil {
			logger.Error("Ошибка получения времени истечения токена из базы: " + err.Error())
			http.Error(w, "Ошибка получения времени истечения токена", http.StatusInternalServerError)
			return
		}

		tokenExpirationRounded := tokenExpiration.Truncate(time.Second)
		claimsExpirationRounded := claims.ExpiresAt.UTC().Truncate(time.Second)

		// Проверка на соответствие времени истечения
		if !tokenExpirationRounded.Equal(claimsExpirationRounded) {
			logger.Info("Токен истек, доступ запрещен.")
			http.Error(w, "Токен истек", http.StatusUnauthorized)
			return
		}

		// Добавление данных пользователя в контекст
		r.Header.Set("userEmail", claims.Email)

		// Передача управления следующему обработчику
		next.ServeHTTP(w, r)
	})
}

// WriteHeader перехватывает статус ответа
func (rw *ResponseWriterWrapper) WriteHeader(code int) {
	rw.StatusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// LoggingMiddleware логирует все запросы, включая статус ответа и ошибки
func LoggingMiddleware(app *internal.App) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			start := time.Now()
			userID := ""

			wrapper := ResponseWriterWrapper{w, 200}

			next.ServeHTTP(&wrapper, r)

			if authHeader := r.Header.Get("Authorization"); authHeader != "" {

				// Проверяем, что заголовок начинается с "Bearer "
				if strings.HasPrefix(authHeader, "Bearer ") {

					// Извлекаем токен
					tokenStr := strings.Split(authHeader, "Bearer ")[1]

					// Декодируем токен и получаем userID
					userId, err := utils.GetUserIDFromToken(tokenStr)
					if err != nil {
						logger.Error("Недействительный токен: " + err.Error())
						http.Error(w, "Недействительный токен", http.StatusUnauthorized)
						return
					}

					// Присваиваем userID
					userID = strconv.Itoa(userId)
				}
			}

			err := app.RequestLogger.Log(
				r.Method,
				r.URL.Path,
				userID,
				r.RemoteAddr,
				r.UserAgent(),
				wrapper.StatusCode,
				time.Since(start),
			)

			if err != nil {
				// Логируем ошибку записи лога, если она произошла
				logger.Error("Ошибка при записи лога запроса:" + err.Error())
			}
		})
	}
}
