// @title Cloud Application API
// @version 1.0
// @host localhost:8081
// @BasePath /
package main

import (
	"Cloud/dataBase"
	_ "Cloud/docs"
	"Cloud/internal"
	"Cloud/logger"
	"Cloud/routes"
	"context"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
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

	// Загрузка переменных
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Ошибка загрузки .env файла")
	}

	// Инициализация логирования
	logger.Logging()

	// Подключение к PostgresSQL
	db := dataBase.ConnectPostgresDB()
	defer db.Close()

	// Подключение к MongoDB
	client := dataBase.ConnectMongoDB()
	defer client.Disconnect(context.Background())

	// Создаем экземпляр App с логгером запросов
	app := &internal.App{
		RequestLogger: logger.NewRequestLogger(client, "Cloud", "logs"),
	}

	// Инициализация маршрутов
	router := routes.InitializeRoutes(db, client, app)

	// Инициализация порта сервера
	portAPI := os.Getenv("SERVER_PORT")

	// Запуск сервера на порту 8081
	go func() {
		wg.Add(1)
		defer wg.Done()
		if err := http.ListenAndServe(":"+portAPI, router); err != nil {
			logger.Error("Failed to start server: " + err.Error())
		}
	}()

	// Логируем запуск сервера
	logger.Info("Server started on port " + portAPI)

	// Ожидаем завершения всех горутин
	wg.Wait()

	// Логируем остановку сервера
	logger.Error("Server has stopped.")
}

// TODO залить
//   Покрыть тестами проект (минимкм 60 %)
