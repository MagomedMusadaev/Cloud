package models

import "time"

// RequestLog представляет собой структуру для логов запросов
type RequestLog struct {
	Method     string        `json:"method"`
	Endpoint   string        `json:"endpoint"`
	UserID     string        `json:"user_id"`
	IP         string        `json:"ip"`
	UserAgent  string        `json:"user_agent"`
	Time       time.Time     `json:"time"`
	StatusCode int           `json:"status_code"` // Статус ответа
	Duration   time.Duration `json:"duration"`    // Время выполнения запроса
}
