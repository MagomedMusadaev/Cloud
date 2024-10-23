package routes

import (
	"Cloud/auth"
	"Cloud/handlers"
	"Cloud/internal"
	"database/sql"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

// InitializeRoutes инициализирует маршруты приложения.
// @Title API Routes
// @Description Настраивает маршруты для операций с пользователями.
func InitializeRoutes(db *sql.DB, client *mongo.Client, app *internal.App) *mux.Router {
	r := mux.NewRouter()

	// Подключаем логирующее middleware
	r.Use(auth.LoggingMiddleware(app))

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)   // Установка статус ответа
		w.Write([]byte("Hello World")) // Запись "Hello World" в тело ответа
	}).Methods("GET")

	// @Summary Создание нового пользователя
	// @Description Создает нового пользователя в системе.
	// @Accept json
	// @Produce json
	// @Param user body models.User true "Пользователь"
	// @Success 201 {string} string "Пользователь успешно создан"
	// @Failure 400 {string} string "Ошибка валидации"
	// @Router /user [post]
	r.HandleFunc("/user", handlers.CreateUser(db)).Methods("POST")

	// @Summary Получение информации о пользователе
	// @Description Получает информацию о пользователе по его уникальному идентификатору.
	// @Produce json
	// @Param id path int true "ID пользователя"
	// @Success 200 {object} models.User "Информация о пользователе"
	// @Failure 400 {string} string "Ошибка при получении пользователя"
	// @Router /user/{id} [get]
	r.HandleFunc("/user/{id}", handlers.GetUser(db)).Methods("GET")

	// @Summary Обновление пользователя
	// @Description Обновляет информацию о пользователе.
	// @Accept json
	// @Produce json
	// @Param id path int true "ID пользователя"
	// @Param user body models.User true "Обновленный пользователь"
	// @Success 204 {string} string "Пользователь успешно обновлен"
	// @Failure 400 {string} string "Ошибка при обновлении пользователя"
	// @Router /user/{id} [put]
	r.HandleFunc("/user/{id}", handlers.UpdateUser(db)).Methods("PUT")

	// @Summary Удаление пользователя
	// @Description Удаляет пользователя из системы по его ID.
	// @Param id path int true "ID пользователя"
	// @Success 204 {string} string "Пользователь успешно удален"
	// @Failure 400 {string} string "Ошибка при удалении пользователя"
	// @Router /user/{id} [delete]
	r.HandleFunc("/user/{id}", handlers.DeleteUser(db)).Methods("DELETE")

	// @Summary Получение всех пользователей
	// @Description Получает список всех пользователей в системе.
	// @Produce json
	// @Success 200 {array} models.User "Список пользователей"
	// @Failure 400 {string} string "Ошибка при получении пользователей"
	// @Router /users [get]
	r.HandleFunc("/users", handlers.GetAllUsers(db)).Methods("GET")

	// @Summary Получение всех логов запросов
	// @Description Получает список всех логов запросов из системы.
	// @Produce application/vnd.openxmlformats-officedocument.spreadsheetml.sheet
	// @Success 200 {file} file "Excel-файл с логами запросов"
	// @Failure 400 {string} string "Ошибка при получении логов"
	// @Router /logs [get]
	r.HandleFunc("/logs", handlers.GetAllRequestLogs(client)).Methods("GET")

	// @Summary Регистрация пользователя
	// @Description Регистрирует нового пользователя в системе.
	// @Accept json
	// @Produce json
	// @Param user body models.User true "Пользователь"
	// @Success 201 {string} string "Пользователь успешно зарегистрирован"
	// @Failure 400 {string} string "Ошибка валидации"
	// @Router /register [post]
	r.HandleFunc("/register", auth.RegisterUser(db)).Methods("POST")

	// @Summary Вход пользователя
	// @Description Позволяет пользователю войти в систему.
	// @Accept json
	// @Produce json
	// @Param user body models.User true "Пользователь"
	// @Success 200 {string} string "Пользователь успешно вошел"
	// @Failure 400 {string} string "Ошибка при входе"
	// @Router /login [post]
	r.HandleFunc("/login", auth.LoginUser(db)).Methods("POST")

	// @Summary Выход пользователя
	// @Description Позволяет пользователю выйти из системы.
	// @Success 200 {string} string "Пользователь успешно вышел"
	// @Failure 400 {string} string "Ошибка при выходе"
	// @Router /logout [post]
	r.HandleFunc("/logout", auth.LogoutHandler(db)).Methods("POST")

	// @Summary Подтверждение электронной почты
	// @Description Подтверждает электронную почту пользователя.
	// @Accept json
	// @Produce json
	// @Param email body string true "Электронная почта пользователя"
	// @Param code body string true "Код подтверждения"
	// @Success 200 {string} string "Электронная почта успешно подтверждена"
	// @Failure 400 {string} string "Ошибка при подтверждении электронной почты"
	// @Router /confirm-email [post]
	r.HandleFunc("/confirm-email", auth.ConfirmEmailHandler(db)).Methods("POST")

	// @Summary Повторная отправка письма с подтверждением
	// @Description Позволяет повторно отправить письмо с подтверждением на электронную почту.
	// @Param email body string true "Электронная почта пользователя"
	// @Success 200 {string} string "Письмо с подтверждением успешно отправлено"
	// @Failure 400 {string} string "Ошибка при повторной отправке"
	// @Router /resend-confirmation [post]
	r.HandleFunc("/resend-confirmation", auth.ResendConfirmationEmailHandler()).Methods("POST")

	// @Summary Обновление access токена
	// @Description Позволяет обновить access токен с использованием refresh токена.
	// @Accept json
	// @Produce json
	// @Success 200 {string} string "Токен успешно обновлен"
	// @Failure 401 {string} string "Недействительный токен"
	// @Router /refresh [post]
	r.HandleFunc("/refresh-token", auth.RefreshTokenHandler).Methods("POST")

	// @Summary Защищенный маршрут
	// @Description Позволяет доступ к защищенному ресурсу только с валидным JWT.
	// @Produce json
	// @Success 200 {string} string "Доступ разрешен"
	// @Failure 401 {string} string "Недействительный токен"
	// @Router /protected [get]
	r.Handle("/protected", auth.JWTMiddleware(db, http.HandlerFunc(ProtectedHandler))).Methods("GET")

	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	return r
}

func ProtectedHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK) // TODO Временная функция для проверки логики работы защищенного маршрута
	w.Write([]byte("Доступ разрешен!"))
}
