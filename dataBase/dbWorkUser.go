package dataBase

import (
	"Cloud/logger"
	"Cloud/models"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
)

// DBCreateUser создает нового пользователя в базе данных.
// @Summary Create a new user
// @Description Adds a new user to the database with the provided details.
// @Accept json
// @Produce json
// @Param user body models.User true "User details"
// @Success 201 {object} models.User
// @Failure 400 {object} ErrorResponse
// @Router /user [post]
func DBCreateUser(db *sql.DB, user *models.User) error {
	query := `INSERT INTO users (name, phone, email, password, from_date_create, from_date_update, is_deleted, is_banned) 
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`

	err := db.QueryRow(query, user.Name, user.Phone, user.Email, user.Password, user.FromDateCreate, user.FromDateUpdate, user.IsDeleted, user.IsBanned).Scan(&user.ID)

	return err
}

// DBGetUser получает пользователя по его ID из базы данных.
// @Summary Get user by ID
// @Description Retrieves a user from the database by their ID.
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} models.User
// @Failure 404 {object} ErrorResponse
// @Router /user/{id} [get]
func DBGetUser(db *sql.DB, userID int) (*models.User, error) {
	var user models.User
	query := `SELECT id, name, phone, email, password, from_date_create, from_date_update, is_deleted, is_banned FROM users WHERE id = $1`

	err := db.QueryRow(query, userID).Scan(&user.ID, &user.Name, &user.Phone, &user.Email, &user.Password, &user.FromDateCreate, &user.FromDateUpdate, &user.IsDeleted, &user.IsBanned)
	if err != nil {
		logger.Error("Failed to retrieve data from the database!" + err.Error())
		return nil, err
	}
	return &user, nil
}

// DBGetAllUsers получает всех пользователей с учетом фильтров, лимита и смещения.
// @Summary Get all users
// @Description Retrieves all users from the database with optional filters, limit, and offset.
// @Accept json
// @Produce json
// @Param filters query string false "Filter users by fields"
// @Param limit query int false "Limit number of users"
// @Param offset query int false "Offset for pagination"
// @Success 200 {array} models.User
// @Failure 500 {object} ErrorResponse
// @Router /users [get]
func DBGetAllUsers(db *sql.DB, filters map[string]string, limit, offset int) ([]*models.User, error) {
	// Базовый SQL-запрос
	query := `SELECT id, name, phone, email, password, from_date_create, from_date_update, is_deleted, is_banned FROM users WHERE TRUE`
	args := []interface{}{}
	counter := 1

	// Условие фильтрации в зависимости от того, какие поля переданы
	for field, value := range filters {
		if value != "" {
			query += fmt.Sprintf(" AND %s ILIKE $%d", field, counter)
			args = append(args, "%"+value+"%") // Добавляем значение в args
			counter++
		}
	}

	// Добавляем постраничность (LIMIT и OFFSET)
	query += fmt.Sprintf(" ORDER BY id LIMIT $%d OFFSET $%d", counter, counter+1)
	args = append(args, limit, offset)

	// Выполняем запрос
	rows, err := db.Query(query, args...)
	if err != nil {
		logger.Error("Failed to retrieve data from the database!" + err.Error())
		return nil, err
	}
	defer rows.Close()

	users := make([]*models.User, 0)
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Name, &user.Phone, &user.Email, &user.Password, &user.FromDateCreate, &user.FromDateUpdate, &user.IsDeleted, &user.IsBanned); err != nil {
			logger.Error("Failed to retrieve data from the database!" + err.Error())
			return nil, err
		}
		users = append(users, &user)
	}

	// Проверка на ошибки после завершения чтения
	if err := rows.Err(); err != nil {
		logger.Error("Failed to retrieve data from the database!" + err.Error())
		return nil, err
	}

	return users, nil
}

// DBUpdateUser обновляет данные о пользователе в базе данных.
// @Summary Update user
// @Description Updates the user details in the database.
// @Accept json
// @Produce json
// @Param user body models.User true "Updated user details"
// @Success 204
// @Failure 400 {object} ErrorResponse
// @Router /user [put]
func DBUpdateUser(db *sql.DB, user *models.User) error {
	query := `UPDATE users SET `
	args := []interface{}{}
	setClauses := []string{}

	// Обработка каждого поля
	if user.Name != "" {
		setClauses = append(setClauses, "name=$"+strconv.Itoa(len(args)+1))
		args = append(args, user.Name)
	}
	if user.Phone != "" {
		setClauses = append(setClauses, "phone=$"+strconv.Itoa(len(args)+1))
		args = append(args, user.Phone)
	}
	if user.Email != "" {
		setClauses = append(setClauses, "email=$"+strconv.Itoa(len(args)+1))
		args = append(args, user.Email)
	}
	if user.Password != "" {
		setClauses = append(setClauses, "password=$"+strconv.Itoa(len(args)+1))
		args = append(args, user.Password)
	}
	if user.FromDateUpdate != "" {
		setClauses = append(setClauses, "from_date_update=$"+strconv.Itoa(len(args)+1))
		args = append(args, user.FromDateUpdate)
	}

	// Добавляем статусы пользователя
	if user.IsDeleted != false { // Предполагается, что false - это "не удален"
		setClauses = append(setClauses, "is_deleted = $"+strconv.Itoa(len(args)+1))
		args = append(args, user.IsDeleted)
	}
	if user.IsBanned != false { // Предполагается, что false - это "не забанен"
		setClauses = append(setClauses, "is_banned = $"+strconv.Itoa(len(args)+1))
		args = append(args, user.IsBanned)
	}

	if len(setClauses) == 0 {
		return nil // Ничего не обновлено
	}

	// Создаём запрос UPDATE
	query += strings.Join(setClauses, ", ") + " WHERE id = $" + strconv.Itoa(len(args)+1)
	args = append(args, user.ID)

	_, err := db.Exec(query, args...)
	return err
}

// DBDeleteUser удаляет пользователя из базы данных по его ID.
// @Summary Delete user
// @Description Deletes a user from the database by their ID.
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 204
// @Failure 404 {object} ErrorResponse
// @Router /user/{id} [delete]
func DBDeleteUser(db *sql.DB, userID int) error {
	query := `DELETE FROM users WHERE id = $1`

	_, err := db.Exec(query, userID)

	return err
}

func FindUserByEmail(db *sql.DB, email string) (*models.User, error) {

	var user models.User
	query := `SELECT id, name, phone, email, password, from_date_create, from_date_update, is_deleted, is_banned FROM users WHERE email = $1 AND is_deleted = false` // разобраться с is_deleted и is_banned

	err := db.QueryRow(query, email).Scan(&user.ID, &user.Name, &user.Phone, &user.Email, &user.Password, &user.FromDateCreate, &user.FromDateUpdate, &user.IsDeleted, &user.IsBanned)

	return &user, err
}

func FindUserByPhone(db *sql.DB, phone string) (*models.User, error) {

	var user models.User
	query := `SELECT id, name, phone, email, password, from_date_create, from_date_update, is_deleted, is_banned FROM users WHERE phone = $1 AND is_deleted = false` // разобраться с is_deleted и is_banned

	err := db.QueryRow(query, phone).Scan(&user.ID, &user.Name, &user.Phone, &user.Email, &user.Password, &user.FromDateCreate, &user.FromDateUpdate, &user.IsDeleted, &user.IsBanned)

	return &user, err
}
