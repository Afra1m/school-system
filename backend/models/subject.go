package models

type Subject struct {
	ID        int    `json:"id" db:"id"`
	Name      string `json:"name" db:"name"`
	TeacherID int    `json:"teacher_id" db:"teacher_id"`
}
