package utils

import "golang.org/x/crypto/bcrypt"

// HashPassword хеширует пароль пользователя.
//
// @Summary Хеширование пароля
// @Description Хеширует пароль пользователя с использованием bcrypt.
// @Param password query string true "Пароль пользователя"
// @Success 200 {string} string "Хешированный пароль"
// @Failure 500 {string} string "Ошибка при хешировании пароля"
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// VerifyPassword проверяет соответствие пароля пользователя и хешированного пароля.
//
// @Summary Проверка пароля
// @Description Сравнивает хешированный пароль с введённым паролем пользователя.
// @Param hashedPassword query string true "Хешированный пароль"
// @Param password query string true "Пароль пользователя"
// @Success 200 {string} string "Пароль совпадает"
// @Failure 401 {string} string "Пароль не совпадает"
// @Failure 500 {string} string "Ошибка при проверке пароля"
func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)) // если err = nil, пароли совпали!
}
