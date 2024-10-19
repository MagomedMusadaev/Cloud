package dataBase

import (
	"Cloud/logger"
	"database/sql"
	_ "github.com/lib/pq"
	"log"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

// ConnectDB устанавливает подключение к базе данных PostgreSQL и возвращает объект DB.
// @Summary Подключение к базе данных PostgreSQL
// @Description Устанавливает соединение с базой данных PostgreSQL с использованием строки подключения.
func ConnectDB() *sql.DB {
	// Строка подключения к базе данных PostgreSQL
	connStr := "host=ep-twilight-breeze-a59htcxo.us-east-2.aws.neon.tech port=5432 user=CloudPostgres_owner password=3EICk0orSfnT dbname=CloudPostgres sslmode=require"

	// Открытие подключения к базе данных
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		// Логирование ошибки при неудачном подключении
		logger.Error("Failed to connect to database!" + err.Error())
		log.Fatal(err) // Завершение программы в случае ошибки
	}

	// Проверка соединения с базой данных (ping)
	if err := db.Ping(); err != nil {
		// Логирование ошибки при неудачном пинге
		logger.Error("Failed to ping database!" + err.Error())
		log.Fatal(err) // Завершение программы в случае ошибки
	}

	// Возврат объекта подключения к базе данных
	return db
}
