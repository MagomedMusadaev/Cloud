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

	// Инициализация переменных
	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	hostDB := os.Getenv("POSTGRES_HOST")
	portDB := os.Getenv("POSTGRES_PORT")
	userDB := os.Getenv("POSTGRES_USER")
	passwordDB := os.Getenv("POSTGRES_PASS")
	dbnameDB := os.Getenv("POSTGRES_NAME")

	// Подключение к базе данных
	db := dataBase.ConnectDB(hostDB, portDB, userDB, passwordDB, dbnameDB)
	defer db.Close()

	// Инициализация маршрутов
	router := routes.InitializeRoutes(db)

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
//   Сделать это всё с токенами AccessToken и RefreshToken
//   Регистрация, логин, логаут(выход с аккаунта), подтверждение почты, отправка повторного кода, рефрештокен(gvt) (google почта для отправки кода подтверждения)
//   При удалении пользователя он не должен удаляться с базы, а переменная IsDeleted должен стать true и надо добавить проверку при авторизации пользователя на переменную IsDeleted, если она равна true, то вход не доступен
//   Сделать хранение лого в базе данных (Mongo DB)
//   Покрыть тестами проект (минимкм 60 %)
