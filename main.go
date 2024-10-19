// @title Cloud Application API
// @version 1.0
// @host localhost:8081
// @BasePath /
package main

import (
	"Cloud/dataBase"
	_ "Cloud/docs"
	"Cloud/logger"
	"Cloud/routes"
	"net/http"
)

// @Summary Основная точка входа приложения.
// @Description Инициализирует логирование, подключается к базе данных, настраивает маршруты и запускает HTTP-сервер.
// @Tags main
// @Success 200 {string} string "Сервер успешно запущен"
// @Failure 500 {string} string "Не удалось запустить сервер"
// @Router / [get]
func main() {
	// Инициализация логирования
	logger.Logging()

	// Подключение к базе данных
	db := dataBase.ConnectDB()
	defer db.Close()

	// Инициализация маршрутов
	router := routes.InitializeRoutes(db)

	// Запуск сервера на порту 8081
	if err := http.ListenAndServe(":8081", router); err != nil {
		logger.Error("Failed to start server: " + err.Error())
	}

	// Логируем остановку сервера
	logger.Error("Server has stopped.")
}
