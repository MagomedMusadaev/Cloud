package routes

import (
	"Cloud/auth"
	"Cloud/handlers"
	"database/sql"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
)

// InitializeRoutes инициализирует маршруты приложения.
// @Title API Routes
// @Description Настраивает маршруты для операций с пользователями.
func InitializeRoutes(db *sql.DB) *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)   // Установите статус ответа
		w.Write([]byte("Hello World")) // Запишите "Hello World" в тело ответа
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

	// Обновление информации о пользователе по ID
	// @Summary Обновление пользователя
	// @Description Обновляет информацию о пользователе.
	// @Accept  json
	// @Produce  json
	// @Param id path int true "ID пользователя"
	// @Param user body models.User true "Обновленный пользователь"
	// @Success 204 {string} string "Пользователь успешно обновлен"
	// @Failure 400 {string} string "Ошибка при обновлении пользователя"
	// @Router /user/{id} [put]
	r.HandleFunc("/user/{id}", handlers.UpdateUser(db)).Methods("PUT")

	// Удаление пользователя по ID
	// @Summary Удаление пользователя
	// @Description Удаляет пользователя из системы по его ID.
	// @Param id path int true "ID пользователя"
	// @Success 204 {string} string "Пользователь успешно удален"
	// @Failure 400 {string} string "Ошибка при удалении пользователя"
	// @Router /user/{id} [delete]
	r.HandleFunc("/user/{id}", handlers.DeleteUser(db)).Methods("DELETE")

	// Получение списка всех пользователей
	// @Summary Получение всех пользователей
	// @Description Получает список всех пользователей в системе.
	// @Produce  json
	// @Success 200 {array} models.User "Список пользователей"
	// @Failure 400 {string} string "Ошибка при получении пользователей"
	// @Router /users [get]
	r.HandleFunc("/users", handlers.GetAllUsers(db)).Methods("GET")

	// Регистрация пользователя (закомментировано)
	// @Summary Регистрация пользователя
	// @Description Регистрирует нового пользователя в системе.
	// @Accept  json
	// @Produce  json
	// @Param user body models.User true "Пользователь"
	// @Success 201 {string} string "Пользователь успешно зарегистрирован"
	// @Failure 400 {string} string "Ошибка валидации"
	// @Router /register [post]
	r.HandleFunc("/register", auth.RegisterUser(db)).Methods("POST")

	// Вход пользователя (закомментировано)
	// @Summary Вход пользователя
	// @Description Позволяет пользователю войти в систему.
	// @Accept  json
	// @Produce  json
	// @Param user body models.User true "Пользователь"
	// @Success 200 {string} string "Пользователь успешно вошел"
	// @Failure 400 {string} string "Ошибка при входе"
	// @Router /login [post]
	r.HandleFunc("/login", auth.LoginUser(db)).Methods("POST")

	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	return r
}

//19.10.2024
// + РЕАЛИЗОВАТЬ ХРАНЕНИЕ ДАННЫХ (ПАРОЛЬ И ИМЯ ПОЛЬЗОВАТЕЛЯ) В ЗАШИФРОВАННОМ ФОРМАТЕ хеширование с солью (переделать столбец password (password VARCHAR(255))
// + Надо сделать чтоб при изменении данных, можно было и несколько параметров отправлять а не все
// + Добавить номер телефона в сущность узер и поменять функции записи данных в базу данных + функции вывода данных
// + минимум символов для пароля 8
// + символов для номер 12 и должен быть вид +79896615609 - 12 символов (+) обязателен
//20.10.2024
//   настроить модуль swagger для генерации документации (создать саму документацию) и включить его в проект
//   SEARCH метод для реализации поиска по всем параметрам (name, email, phone)
//   Релизовать методы регистрации и авторизации пользователя с помощью email и password, или телефон и пароль
//   Регистрация, логин, логаут(выход с аккаунта), подтверждение почты, отправка повторного кода, рефрештокен(gvt) (google почта для отправки кода подтверждения)
