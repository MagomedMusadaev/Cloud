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
	"sync"
)

// @Summary Основная точка входа приложения.
// @Description Инициализирует логирование, подключается к базе данных, настраивает маршруты и запускает HTTP-сервер.
// @Tags main
// @Success 200 {string} string "Сервер успешно запущен"
// @Failure 500 {string} string "Не удалось запустить сервер"
// @Router / [get]
func main() {

	wg := sync.WaitGroup{}

	// Инициализация логирования
	logger.Logging()

	// Подключение к базе данных
	db := dataBase.ConnectDB()
	defer db.Close()

	// Инициализация маршрутов
	router := routes.InitializeRoutes(db)

	// Запуск сервера на порту 8081
	go func() {
		wg.Add(1)
		defer wg.Done()
		if err := http.ListenAndServe(":8081", router); err != nil {
			logger.Error("Failed to start server: " + err.Error())
		}
	}()

	// Логируем запуск сервера
	logger.Info("Server started on port 8081")

	// Ожидаем завершения всех горутин
	wg.Wait()

	// Логируем остановку сервера
	logger.Error("Server has stopped.")
}
