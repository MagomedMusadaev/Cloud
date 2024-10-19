package utils

import (
	"Cloud/models"
	"errors"
	"regexp"
)

// ValidateUserForCreate проверяет корректность данных пользователя при создании.
//
// @Summary Валидация данных пользователя для создания
// @Description Проверяет, что имя пользователя, номер телефона, пароль и email заполнены и соответствуют требованиям.
// @Param user body models.User true "Данные пользователя"
// @Success 200 {string} string "Данные пользователя валидны"
// @Failure 400 {string} string "Ошибка валидации данных"
// @Router /user [post]
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

// ValidateUserForUpdate проверяет корректность данных пользователя при обновлении.
//
// @Summary Валидация данных пользователя для обновления
// @Description Проверяет, что имя пользователя, номер телефона, пароль и email соответствуют требованиям.
// @Param user body models.User true "Данные пользователя"
// @Success 200 {string} string "Данные пользователя валидны"
// @Failure 400 {string} string "Ошибка валидации данных"
// @Router /user/{id} [put]
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
		if len(user.Phone) != 12 {
			return errors.New("Phone length is invalid!")
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

// isValidUsername проверяет имя пользователя на допустимые символы.
//
// @Summary Проверка имени пользователя
// @Description Проверяет имя пользователя на допустимые символы: латинские буквы, цифры, дефисы и подчеркивания.
// @Param name query string true "Имя пользователя"
// @Success 200 {boolean} bool "Имя пользователя корректно"
// @Failure 400 {string} string "Имя пользователя некорректно"
func isValidUsername(name string) bool {
	re := regexp.MustCompile(`^[\p{L}\d-_]+$`)
	return re.MatchString(name)
}

// isValidPhone проверяет номер телефона на допустимые символы.
//
// @Summary Проверка номера телефона
// @Description Проверяет номер телефона на допустимые символы: плюс и цифры.
// @Param phone query string true "Номер телефона"
// @Success 200 {boolean} bool "Номер телефона корректен"
// @Failure 400 {string} string "Номер телефона некорректен"
func isValidPhone(phone string) bool {
	re := regexp.MustCompile(`^\+\d+$`)
	return re.MatchString(phone)
}

// isValidPassword проверяет пароль на допустимые символы.
//
// @Summary Проверка пароля
// @Description Проверяет пароль на допустимые символы: латинские буквы, цифры и специальные символы.
// @Param password query string true "Пароль"
// @Success 200 {boolean} bool "Пароль корректен"
// @Failure 400 {string} string "Пароль некорректен"
func isValidPassword(password string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9\W_]+$`)
	return re.MatchString(password)
}

// isValidEmail проверяет email на корректность.
//
// @Summary Проверка email
// @Description Проверяет email на соответствие стандартному формату.
// @Param email query string true "Email"
// @Success 200 {boolean} bool "Email корректен"
// @Failure 400 {string} string "Email некорректен"
func isValidEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}
