package auth

import (
	"database/sql"
	"net/http"
)

// Вход пользователя на свой аккаунт
func LoginUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Логика аутентификации пользователя
	}
}
