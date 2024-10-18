package utils

import (
	"Cloud/models"
	"errors"
	"regexp"
)

func ValidateUserForCreate(user models.User) error {

	if user.Name == "" || !isValidUsername(user.Name) {
		return errors.New("Username is empty!")
	}

	if user.Phone == "" || !isValidPhone(user.Phone) {
		return errors.New("Phone is empty!")
	}

	if user.Password == "" || !isValidPassword(user.Password) {
		return errors.New("Password cannot be empty!")
	}
	if len(user.Password) < 8 {
		return errors.New("Password is too short!")
	}

	if user.Email == "" || !isValidEmail(user.Email) {
		return errors.New("Email is invalid!")
	}

	return nil
}

func ValidateUserForUpdate(user models.User) error {
	if user.Name != "" {
		if !isValidUsername(user.Name) {
			return errors.New("Name is invalid!")
		}
	}

	if user.Phone != "" {
		if !isValidPhone(user.Phone) {
			return errors.New("Phone is invalid!")
		}
	}

	if user.Password != "" {
		if !isValidPassword(user.Password) {
			return errors.New("Password is invalid!")
		}
	}

	if user.Email != "" {
		if !isValidEmail(user.Email) {
			return errors.New("Email is invalid!")
		}
	}

	return nil
}

// Проверка имени пользователя на правильность символов
func isValidUsername(name string) bool {
	re := regexp.MustCompile(`^[\p{L}\d-_]+$`)
	return re.MatchString(name)
}

// Проверка номера телефона пользователя на правильность символов
func isValidPhone(phone string) bool {
	re := regexp.MustCompile(`^\d{3}-\d{3}-\d{3}$`) // 123-456-678 (надо будет потом поменяить на +7 и т.д.
	return re.MatchString(phone)
}

// Проверка пароля пользователя на правильность символов
func isValidPassword(password string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9\W_]+$`)
	return re.MatchString(password)
}

// Проверка email пользователя на правильность символов
func isValidEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}
