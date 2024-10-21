package auth

import (
	"database/sql"
	"net/http"
)

// Регистрация пользователя (не админ)
func RegisterUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Логика регистрации пользователя
	}
}
