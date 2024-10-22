package models

import "github.com/golang-jwt/jwt/v4"

// Claims — это кастомная структура для JWT с дополнительным полем Email
type Claims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}
