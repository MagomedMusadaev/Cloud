package dataBase

import (
	"Cloud/logger"
	"context"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"time"
)

// ErrorResponse представляет стандартный ответ с сообщением об ошибке
type ErrorResponse struct {
	Message string `json:"message"`
}

// ConnectPostgresDB устанавливает подключение к базе данных PostgreSQL и возвращает объект DB.
// @Summary Подключение к базе данных PostgreSQL
// @Description Устанавливает соединение с базой данных PostgreSQL с использованием строки подключения.
// @Success 200 {string} string "Успешное подключение к PostgreSQL"
// @Failure 500 {object} ErrorResponse "Ошибка подключения к PostgreSQL"
// @Router /connect-postgres [get]
func ConnectPostgresDB() *sql.DB {

	hostPostgresDB := os.Getenv("POSTGRES_HOST")
	portPostgresDB := os.Getenv("POSTGRES_PORT")
	userPostgresDB := os.Getenv("POSTGRES_USER")
	passwordPostgresDB := os.Getenv("POSTGRES_PASS")
	dbnamePostgresDB := os.Getenv("POSTGRES_NAME")

	// Строка подключения к базе данных PostgreSQL
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=require",
		hostPostgresDB, portPostgresDB, userPostgresDB, passwordPostgresDB, dbnamePostgresDB)

	// Открытие подключения к базе данных
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		// Логирование ошибки при неудачном подключении
		logger.Error("Failed to connect to PostgresDB!" + err.Error())
		log.Fatal(err) // Завершение программы в случае ошибки
	}

	// Проверка соединения с базой данных (ping)
	if err := db.Ping(); err != nil {
		// Логирование ошибки при неудачном пинге
		logger.Error("Failed to ping PostgresDB!" + err.Error())
		log.Fatal(err) // Завершение программы в случае ошибки
	}

	logger.Info("Успешно подключено к PostgresDB")

	// Возврат объекта подключения к базе данных
	return db
}

// ConnectMongoDB устанавливает подключение к базе данных MongoDB и возвращает объект клиента.
// @Summary Подключение к базе данных MongoDB
// @Description Устанавливает соединение с MongoDB с использованием переменных окружения для формирования строки подключения.
// @Success 200 {string} string "Успешное подключение к MongoDB"
// @Failure 500 {object} ErrorResponse "Ошибка подключения к MongoDB"
// @Router /connect-mongo [get]
func ConnectMongoDB() *mongo.Client {
	// Чтение переменных окружения
	hostMongoDB := os.Getenv("MONGO_HOST")
	userMongoDB := os.Getenv("MONGO_USER")
	passwordMongoDB := os.Getenv("MONGO_PASS")

	// Формирование строки подключения
	mongoURI := fmt.Sprintf("mongodb+srv://%s:%s@%s", userMongoDB, passwordMongoDB, hostMongoDB)

	// Настройка клиента
	clientOptions := options.Client().ApplyURI(mongoURI)

	// Контекст с тайм-аутом
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Подключение к MongoDB с созданием клиента
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal("Ошибка подключения к MongoDB:", err)
	}

	// Проверка соединения
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Не удалось подключиться к MongoDB:", err)
	}

	logger.Info("Успешно подключено к MongoDB")

	return client
}
