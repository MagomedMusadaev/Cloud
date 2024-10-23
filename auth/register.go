package auth

import (
	"Cloud/email"
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

		// Генерация кода подтверждения
		confirmationCode := utils.GenRandCode()                         // создайте эту функцию для генерации кода
		err = email.SendConfirmationEmail(user.Email, confirmationCode) // отправка кода на почту
		if err != nil {
			logger.Error("Ошибка отправки письма: " + err.Error())
			http.Error(w, "Ошибка отправки письма", http.StatusInternalServerError)
			return
		}

		// Сохраняем код, данные пользователя и время создания в TemporaryStore на 1 час
		models.TemporaryStore[user.Email] = models.ConfirmationData{
			Code:      confirmationCode,
			User:      user,
			CreatedAt: time.Now(),
		}

		// Ответ клиенту
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Код подтверждения отправлен на " + user.Email})
	}
}
