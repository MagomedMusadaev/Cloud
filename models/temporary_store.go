package models

import "time"

// Структура для хранения кода подтверждения, данных пользователя и времени создания
type ConfirmationData struct {
	Code      string
	User      User
	CreatedAt time.Time
}

// Временное хранилище для подтверждения email
var TemporaryStore = make(map[string]ConfirmationData)
