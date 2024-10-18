package dataBase

import (
	"Cloud/logger"
	"Cloud/models"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
)

// Создание пользователя в db
func DBCreateUser(db *sql.DB, user *models.User) error {

	query := `INSERT INTO users (name, phone, email, password, from_date_create, from_date_update, is_deleted, is_banned) 
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`

	err := db.QueryRow(query, user.Name, user.Phone, user.Email, user.Password, user.FromDateCreate, user.FromDateUpdate, user.IsDeleted, user.IsBanned).Scan(&user.ID)

	return err
}

// Получение пользователя из db
func DBGetUser(db *sql.DB, userID int) (*models.User, error) {

	var user models.User
	query := `SELECT id, name, phone, email,password, from_date_create, from_date_update, is_deleted, is_banned FROM users WHERE id = $1`

	err := db.QueryRow(query, userID).Scan(&user.ID, &user.Name, &user.Phone, &user.Email, &user.Password, &user.FromDateCreate, &user.FromDateUpdate, &user.IsDeleted, &user.IsBanned)
	if err != nil {
		logger.Error("Failed to retrieve data from the database!" + err.Error())
		return nil, err
	}
	return &user, nil

}

// Получение пользователей, с учётом лимита на страницу и смещения страницы, из db
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

// Изменение данных о пользователе в db
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

// Удаление пользователя из db
func DBDeleteUser(db *sql.DB, userID int) error {
	query := `DELETE FROM users WHERE id = $1`

	_, err := db.Exec(query, userID)

	return err
}
