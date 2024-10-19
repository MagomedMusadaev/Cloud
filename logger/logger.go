package logger

import (
	"fmt"
	"io"
	"log"
	"os"
)

var (
	infoLog    *log.Logger // Логгер для информационных сообщений
	warningLog *log.Logger // Логгер для предупреждающих сообщений
	errorLog   *log.Logger // Логгер для сообщений об ошибках
	LogFile    *os.File    // Файл для записи логов
)

// Logging инициализирует логгеры и открывает файл для записи логов.
// @Summary Инициализация логгирования
// @Description Открывает файл для записи логов и инициализирует логгеры для разных уровней
// @Tags logger
// @Success 200 "Logging initialized successfully"
// @Failure 500 "Failed to initialize logging"
// @Router /logger/init [post]
func Logging() {
	// Открываем файл для записи логов
	logFile, err := os.OpenFile("Cloud.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Не удалось открыть файл для логов!", err)
	}

	// Создаем MultiWriter для записи в файл и в консоль
	multiWriter := io.MultiWriter(os.Stdout, logFile)

	// Инициализируем логгеры для разных уровней
	infoLog = log.New(multiWriter, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	warningLog = log.New(multiWriter, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	errorLog = log.New(multiWriter, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}

// Info записывает информационное сообщение в лог.
// @Summary Запись информационного сообщения
// @Description Записывает информационное сообщение в лог
// @Tags logger
// @Param msg query string true "Сообщение для логирования"
// @Success 200 "Message logged successfully"
// @Failure 500 "Failed to log message"
// @Router /logger/info [post]
func Info(msg string) {
	if infoLog != nil {
		infoLog.Println(msg)
	} else {
		fmt.Println("INFO LOGGING NOT INITIALIZED")
	}
}

// Warning записывает предупреждающее сообщение в лог.
// @Summary Запись предупреждающего сообщения
// @Description Записывает предупреждающее сообщение в лог
// @Tags logger
// @Param msg query string true "Сообщение для логирования"
// @Success 200 "Message logged successfully"
// @Failure 500 "Failed to log message"
// @Router /logger/warning [post]
func Warning(msg string) {
	if warningLog != nil {
		warningLog.Println(msg)
	} else {
		fmt.Println("WARNING LOGGING NOT INITIALIZED")
	}
}

// Error записывает сообщение об ошибке в лог.
// @Summary Запись сообщения об ошибке
// @Description Записывает сообщение об ошибке в лог
// @Tags logger
// @Param msg query string true "Сообщение для логирования"
// @Success 200 "Message logged successfully"
// @Failure 500 "Failed to log message"
// @Router /logger/error [post]
func Error(msg string) {
	if errorLog != nil {
		errorLog.Println(msg)
	} else {
		fmt.Println("ERROR LOGGING NOT INITIALIZED")
	}
}
