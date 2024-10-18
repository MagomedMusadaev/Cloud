package security

import "golang.org/x/crypto/bcrypt"

// Функция хенирования пароля пользователя
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// Проверка на соответствие пароля пользователя и пароля в db
func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)) // если err = nil, пароли совпали!
}
