package handlers

import (
	"Cloud/dataBase"
	"Cloud/logger"
	"Cloud/models"
	"Cloud/security"
	"Cloud/utils"
	"database/sql"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"time"
)

func CreateUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var user models.User

		//Декодирование JSON из тела запроса
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			logger.Error("Failed to connect to database: " + err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		//Валидация данных пользователя
		if err := utils.ValidateUserForCreate(user); err != nil {
			logger.Error("User validation failed!" + err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		//Установка даты создания и обновления
		user.FromDateCreate = time.Now().Format(time.RFC3339)
		user.FromDateUpdate = user.FromDateCreate

		// Хеширование пароля перед сохранением
		user.Password, err = security.HashPassword(user.Password)
		if err != nil {
			logger.Error("Failed to hash password: " + err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		//Запрос к базе данных
		err = dataBase.DBCreateUser(db, &user)
		if err != nil {
			logger.Error("Failed to create user: " + err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Успешный ответ
		w.WriteHeader(http.StatusCreated)                                                  // устанавливается статус ответа 201 Created
		json.NewEncoder(w).Encode("User " + user.Name + " has been successfully created!") // сериализует объект "user" в JSON и отправляет в "w"(ответ)
	}
}

func GetUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		params := mux.Vars(r)
		userID, err := strconv.Atoi(params["id"])
		if err != nil {
			logger.Error("Failed to parse id in 'GetUser': " + err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		//Запрос к базе данных
		user, err := dataBase.DBGetUser(db, userID)
		if err != nil {
			logger.Error("Failed to get user: " + err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Успешный ответ
		json.NewEncoder(w).Encode(user)
	}
}

func GetAllUsers(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		pageStr := r.URL.Query().Get("page")
		limitStr := r.URL.Query().Get("limit")
		name := r.URL.Query().Get("name")
		email := r.URL.Query().Get("email")
		phone := r.URL.Query().Get("phone")

		// Устанавливаем значения по умолчанию
		page := 1
		limit := 10

		// Преобразуем параметры постраничной навигации
		if pageStr != "" {
			p, err := strconv.Atoi(pageStr)
			if err == nil && p > 0 {
				page = p
			}
		}
		if limitStr != "" {
			l, err := strconv.Atoi(limitStr)
			if err == nil && l > 0 {
				limit = l
			}
		}

		//Для первой страницы: offset = (1 - 1) * 5 = 0, что означает, что записи будут начинаться с первого пользователя.
		//Для первой страницы: offset = (1 - 1) * 5 = 0, что означает, что записи будут начинаться с первого пользователя.
		offset := (page - 1) * limit // Расчёт смещения для SQL-запроса на основе номера страницы и лимита

		// Собираем фильтры
		filter := map[string]string{
			"name":  name,
			"email": email,
			"phone": phone,
		}

		// Получаем пользователей из базы данных с учётом фильтров и постраничности
		users, err := dataBase.DBGetAllUsers(db, filter, limit, offset)
		if err != nil {
			logger.Error("Failed to get all users: " + err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Возвращаем пользователей в формате JSON
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(users)
	}
}

func UpdateUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		params := mux.Vars(r)
		userID, err := strconv.Atoi(params["id"])
		if err != nil {
			logger.Error("Failed to parse id in 'UpdateUser': " + err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var user models.User
		err = json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			logger.Error("Failed to decode user: " + err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		user.ID = userID

		user.FromDateUpdate = time.Now().Format(time.RFC3339)

		if err := utils.ValidateUserForUpdate(user); err != nil {
			logger.Error("User validation failed!" + err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Хеширование пароля если он был изменён
		if user.Password != "" {
			user.Password, err = security.HashPassword(user.Password)
			if err != nil {
				logger.Error("Failed to hash password: " + err.Error())
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		// Обновление данных пользователя в базе
		err = dataBase.DBUpdateUser(db, &user)
		if err != nil {
			logger.Error("Failed to update user: " + err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Ответ без тела (204 No Content)
		w.WriteHeader(http.StatusNoContent)
	}
}

func DeleteUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		params := mux.Vars(r)
		userID, err := strconv.Atoi(params["id"])
		if err != nil {
			logger.Error("Failed to parse id in 'DeleteUser': " + err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = dataBase.DBDeleteUser(db, userID)
		if err != nil {
			logger.Error("Failed to delete user: " + err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
