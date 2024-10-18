package models

type User struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	Phone          string `json:"phone"`
	Email          string `json:"email"`
	Password       string `json:"password"`
	FromDateCreate string `json:"fromDateCreate"`
	FromDateUpdate string `json:"fromDateUpdate"`
	IsDeleted      bool   `json:"isDeleted"`
	IsBanned       bool   `json:"isBanned"`
}
