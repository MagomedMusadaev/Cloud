package utils

import (
	"Cloud/models"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
)

// GetUserIDFromToken извлекает userID из токена
func GetUserIDFromToken(token string) (int, error) {
	// Разделение токена на части
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return 0, fmt.Errorf("недействительный токен")
	}

	// Декодирование полезной нагрузки
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return 0, fmt.Errorf("ошибка декодирования полезной нагрузки: %v", err)
	}

	// Распаковка полезной нагрузки в структуру Claims
	var claims models.Claims
	if err := json.Unmarshal(payload, &claims); err != nil {
		return 0, fmt.Errorf("ошибка распаковки полезной нагрузки: %v", err)
	}

	return claims.UserID, nil // Возвращаем userID
}
