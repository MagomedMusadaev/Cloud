package routes

import (
	"Cloud/handlers"
	"database/sql"
	"github.com/gorilla/mux"
)

// InitializeRoutes инициализирует маршруты приложения.
func InitializeRoutes(db *sql.DB) *mux.Router {
	r := mux.NewRouter()

	// Определения маршрутов для CRUD операций с пользователями
	r.HandleFunc("/user", handlers.CreateUser(db)).Methods("POST")
	r.HandleFunc("/user/{id}", handlers.GetUser(db)).Methods("GET")
	r.HandleFunc("/user/{id}", handlers.UpdateUser(db)).Methods("PUT")
	r.HandleFunc("/user/{id}", handlers.DeleteUser(db)).Methods("DELETE")
	r.HandleFunc("/users", handlers.GetAllUsers(db)).Methods("GET")
	//r.HandleFunc("/register")
	//r.HandleFunc("/login")

	return r
}

// + РЕАЛИЗОВАТЬ ХРАНЕНИЕ ДАННЫХ (ПАРОЛЬ И ИМЯ ПОЛЬЗОВАТЕЛЯ) В ЗАШИФРОВАННОМ ФОРМАТЕ хеширование с солью (переделать столбец password (password VARCHAR(255))
// + Надо сделать чтоб при изменении данных, можно было и несколько параметров отправлять а не все
// + Добавить номер телефона в сущность узер и поменять функции записи данных в базу данных + функции вывода данных
// минимум символов для пароля 8 номер 11
