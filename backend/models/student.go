package models

type Student struct {
	ID        int    `json:"id" db:"id"`
	FullName  string `json:"full_name" db:"full_name"`
	ClassName string `json:"class_name" db:"class_name"`
	UserID    int    `json:"user_id" db:"user_id"`
}
