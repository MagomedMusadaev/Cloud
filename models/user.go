package models

// User представляет пользователя в системе.
// @Description Модель пользователя с основными полями.
// @Title User
// @Required
type User struct {
	// @Description Уникальный идентификатор пользователя
	// @Example 1
	ID int `json:"id"`

	// @Description Имя пользователя
	// @Example "John Doe"
	Name string `json:"name"`

	// @Description Номер телефона пользователя
	// @Example "+1234567890"
	Phone string `json:"phone"`

	// @Description Электронная почта пользователя
	// @Example "johndoe@example.com"
	Email string `json:"email"`

	// @Description Пароль пользователя (сохраняется в зашифрованном виде)
	// @Example "password123"
	Password string `json:"password"`

	// @Description Дата и время создания пользователя
	// @Example "2024-10-19T12:00:00Z"
	FromDateCreate string `json:"fromDateCreate"`

	// @Description Дата и время последнего обновления пользователя
	// @Example "2024-10-19T12:00:00Z"
	FromDateUpdate string `json:"fromDateUpdate"`

	// @Description Флаг, указывающий, удален ли пользователь
	// @Example false
	IsDeleted bool `json:"isDeleted"`

	// @Description Флаг, указывающий, заблокирован ли пользователь
	// @Example false
	IsBanned bool `json:"isBanned"`
}
