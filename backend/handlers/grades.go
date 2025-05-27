package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"school-system/backend/database"
	"school-system/backend/models"

	"github.com/gorilla/mux"
)

func GetGrades(w http.ResponseWriter, r *http.Request) {
	log.Printf("Получен запрос на получение списка оценок")
	row, err := database.DB.Query(`SELECT id, student_id, subject_id, grade, quarter FROM grades `)
	if err != nil {
		log.Printf("Ошибка при получении данных оценок: %v", err)
		http.Error(w, "Ошибка при получении данных", http.StatusInternalServerError)
		return
	}
	defer row.Close()

	var grades []models.Grade
	for row.Next() {
		var g models.Grade
		err := row.Scan(&g.ID, &g.StudentID, &g.SubjectID, &g.Grade, &g.Quarter)
		if err != nil {
			log.Printf("Ошибка при сканировании данных оценки: %v", err)
			http.Error(w, "Ошибка при получении данных", http.StatusInternalServerError)
			return
		}
		grades = append(grades, g)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(grades)
}

func CreateGrade(w http.ResponseWriter, r *http.Request) {
	log.Printf("Получен запрос на создание новой оценки")
	var grade models.Grade
	if err := json.NewDecoder(r.Body).Decode(&grade); err != nil {
		log.Printf("Ошибка при чтении данных оценки: %v", err)
		http.Error(w, "Невозможно прочитать данные", http.StatusBadRequest)
		return
	}

	_, err := database.DB.Exec(`INSERT INTO grades (student_id, subject_id, grade, quarter) VALUES ($1, $2, $3, $4)`,
		grade.StudentID, grade.SubjectID, grade.Grade, grade.Quarter)
	if err != nil {
		log.Printf("Ошибка при добавлении оценки в базу данных: %v", err)
		http.Error(w, "Ошибка при добавлении оценки", http.StatusInternalServerError)
		return
	}

	log.Printf("Успешно создана новая оценка для студента %d по предмету %d", grade.StudentID, grade.SubjectID)
	w.WriteHeader(http.StatusCreated)
}

func UpdateGrade(w http.ResponseWriter, r *http.Request) {
	log.Printf("Получен запрос на обновление оценки")
	var g models.Grade
	if err := json.NewDecoder(r.Body).Decode(&g); err != nil {
		log.Printf("Ошибка при чтении данных оценки: %v", err)
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}

	_, err := database.DB.Exec(`UPDATE grades SET student_id=$1, subject_id=$2, grade=$3, quarter=$4 WHERE id=$5`,
		g.StudentID, g.SubjectID, g.Grade, g.Quarter, g.ID)
	if err != nil {
		log.Printf("Ошибка при обновлении оценки с ID %d: %v", g.ID, err)
		http.Error(w, "Ошибка при обновлении", http.StatusInternalServerError)
		return
	}

	log.Printf("Успешно обновлена оценка с ID: %d", g.ID)
	w.WriteHeader(http.StatusOK)
}

func DeleteGrade(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	log.Printf("Получен запрос на удаление оценки с ID: %s", id)

	// Получаем информацию об оценке перед удалением
	var grade models.Grade
	err := database.DB.Get(&grade, `SELECT * FROM grades WHERE id=$1`, id)
	if err != nil {
		log.Printf("Ошибка при получении информации об оценке с ID %s: %v", id, err)
		http.Error(w, "Ошибка при удалении", http.StatusInternalServerError)
		return
	}

	// Удаляем оценку
	_, err = database.DB.Exec(`DELETE FROM grades WHERE id=$1`, id)
	if err != nil {
		log.Printf("Ошибка при удалении оценки с ID %s: %v", id, err)
		http.Error(w, "Ошибка при удалении", http.StatusInternalServerError)
		return
	}

	log.Printf("Успешно удалена оценка: ID=%s, Студент=%d, Предмет=%d, Оценка=%d, Четверть=%d",
		id, grade.StudentID, grade.SubjectID, grade.Grade, grade.Quarter)
	w.WriteHeader(http.StatusOK)
}

func GetStudentGrades(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	studentID := vars["id"]
	log.Printf("Получен запрос на получение оценок студента с ID: %s", studentID)

	query := `
		SELECT g.*, s.name as subject_name
		FROM grades g
		JOIN subjects s ON g.subject_id = s.id
		WHERE g.student_id = $1
		ORDER BY s.name, g.quarter
	`

	rows, err := database.DB.Queryx(query, studentID)
	if err != nil {
		log.Printf("Ошибка при получении оценок студента %s: %v", studentID, err)
		http.Error(w, "Ошибка при получении оценок", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type GradeWithSubject struct {
		models.Grade
		SubjectName string `db:"subject_name" json:"subject_name"`
	}

	var grades []GradeWithSubject
	for rows.Next() {
		var grade GradeWithSubject
		if err := rows.StructScan(&grade); err != nil {
			log.Printf("Ошибка при обработке данных оценки: %v", err)
			http.Error(w, "Ошибка при обработке данных", http.StatusInternalServerError)
			return
		}
		grades = append(grades, grade)
	}

	log.Printf("Успешно получено %d оценок для студента с ID: %s", len(grades), studentID)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(grades)
}
