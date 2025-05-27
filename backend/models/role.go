package models

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"` // здесь будет храниться хэш пароля
	Role     string `json:"role"`     // student, teacher, deputy
}
