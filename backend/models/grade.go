package models

type Grade struct {
	ID        int `json:"id" db:"id"`
	StudentID int `json:"student_id" db:"student_id"`
	SubjectID int `json:"subject_id" db:"subject_id"`
	Grade     int `json:"grade" db:"grade"`
	Quarter   int `json:"quarter" db:"quarter"`
}
