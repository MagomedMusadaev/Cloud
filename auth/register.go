package auth

import (
	"Cloud/dataBase"
	"Cloud/logger"
	"Cloud/models"
	"Cloud/utils"
	"database/sql"
	"encoding/json"
	"net/http"
	"time"
)

// Регистрация пользователя (не админ)
func RegisterUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var user models.User

		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			logger.Error("Ошибка декодирования JSON: " + err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		//Валидация данных пользователя
		if err := utils.ValidateUserForCreate(user); err != nil {
			logger.Error("User validation failed!" + err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Установка даты и создания обновления
		user.FromDateCreate = time.Now().Format(time.RFC3339)
		user.FromDateUpdate = user.FromDateCreate

		// Хеширование пароля перед сохранением
		user.Password, err = utils.HashPassword(user.Password)
		if err != nil {
			logger.Error("Failed to hash password: " + err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		//Запрос к базе данных
		err = dataBase.DBCreateUser(db, &user)
		if err != nil {
			logger.Error("Failed to create user: " + err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Успешный ответ
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode("User " + user.Name + " has been successfully created!")
	}
}
