package handlers

import (
	"Cloud/dataBase"
	"Cloud/logger"
	"Cloud/models"
	"Cloud/utils"
	"database/sql"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"time"
)

// CreateUser handles user creation.
// @Summary Create a new user
// @Description Create a new user with the provided details
// @Tags users
// @Accept json
// @Produce json
// @Param user body models.User true "User data"
// @Success 201 {string} string "User created successfully"
// @Failure 400 {string} string "Invalid request format"
// @Failure 500 {string} string "Internal server error"
// @Router /users [post]
func CreateUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User

		//Декодирование JSON из тела запроса
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			logger.Error("Ошибка декодирования JSON: " + err.Error())
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
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Успешный ответ
		w.WriteHeader(http.StatusCreated)                                                  // устанавливается статус ответа 201 Created
		json.NewEncoder(w).Encode("User " + user.Name + " has been successfully created!") // сериализует объект "user" в JSON и отправляет в "w"(ответ)
	}
}

// GetUser retrieves a user by ID.
// @Summary Get a user by ID
// @Description Retrieve a user using their ID
// @Tags users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} models.User "User data"
// @Failure 400 {string} string "Invalid ID"
// @Failure 404 {string} string "User not found"
// @Router /users/{id} [get]
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

// GetAllUsers получает список пользователей с постраничной навигацией и фильтрами.
// @Summary Получить всех пользователей
// @Description Получение списка пользователей с возможностью фильтрации по полям
// @Tags users
// @Produce json
// @Param page query int false "Номер страницы"
// @Param limit query int false "Количество пользователей на странице"
// @Param name query string false "Фильтр по имени"
// @Param email query string false "Фильтр по email"
// @Param phone query string false "Фильтр по телефону"
// @Success 200 {array} models.User "Список пользователей"
// @Failure 400 {string} string "Некорректный запрос"
// @Router /users [get]
func GetAllUsers(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pageStr := r.URL.Query().Get("page")
		limitStr := r.URL.Query().Get("limit")
		name := r.URL.Query().Get("name")
		email := r.URL.Query().Get("email")
		phone := r.URL.Query().Get("phone")
		searchWord := r.URL.Query().Get("search")

		if name == "" {
			name = searchWord
		}
		if email == "" {
			email = searchWord
		}
		if phone == "" {
			phone = searchWord
		}

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

// UpdateUser updates a user's information.
// @Summary Update a user's information
// @Description Update a user's details by their ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param user body models.User true "User data"
// @Success 204 "User updated successfully"
// @Failure 400 {string} string "Invalid request"
// @Failure 404 {string} string "User not found"
// @Router /users/{id} [put]
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
			user.Password, err = utils.HashPassword(user.Password)
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

// DeleteUser removes a user by ID.
// @Summary Delete a user by ID
// @Description Remove a user using their ID
// @Tags users
// @Param id path int true "User ID"
// @Success 204 "User deleted successfully"
// @Failure 400 {string} string "Invalid ID"
// @Failure 404 {string} string "User not found"
// @Router /users/{id} [delete]
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

func GetLogs() {

}
