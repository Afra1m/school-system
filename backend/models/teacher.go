package models

type Teacher struct {
	ID         int    `json:"id" db:"id"`
	FullName   string `json:"full_name" db:"full_name"`
	SubjectID  int    `json:"subject_id" db:"subject_id"`
	RoomNumber string `json:"room_number" db:"room_number"`
	UserID     int    `json:"user_id,omitempty" db:"user_id"`
}
