package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"school-system/backend/database"
)

// GetStudentsCount возвращает общее количество учеников
func GetStudentsCount(w http.ResponseWriter, r *http.Request) {
	var count int64
	err := database.DB.Get(&count, "SELECT COUNT(*) FROM students")
	if err != nil {
		log.Printf("Ошибка при получении количества учеников: %v", err)
		http.Error(w, "Ошибка при получении количества учеников: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]int64{"count": count})
}

// GetTeachersCount возвращает общее количество учителей
func GetTeachersCount(w http.ResponseWriter, r *http.Request) {
	var count int64
	err := database.DB.Get(&count, "SELECT COUNT(*) FROM teachers")
	if err != nil {
		log.Printf("Ошибка при получении количества учителей: %v", err)
		http.Error(w, "Ошибка при получении количества учителей: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]int64{"count": count})
}

// GetAverageGrade возвращает средний балл по всем оценкам
func GetAverageGrade(w http.ResponseWriter, r *http.Request) {
	var avg sql.NullFloat64
	query := `
		SELECT COALESCE(AVG(grade), 0) 
		FROM grades 
		WHERE grade IS NOT NULL
	`
	err := database.DB.Get(&avg, query)
	if err != nil {
		log.Printf("Ошибка при получении среднего балла: %v", err)
		http.Error(w, "Ошибка при получении среднего балла: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]float64{"average": avg.Float64})
}

// GetClassPerformance возвращает средние оценки по классам
func GetClassPerformance(w http.ResponseWriter, r *http.Request) {
	type ClassPerformance struct {
		Name  string  `json:"name" db:"class_name"`
		Value float64 `json:"value" db:"average"`
	}

	query := `
		WITH class_quarter_averages AS (
			SELECT 
				s.class_name,
				g.quarter,
				AVG(g.grade) as quarter_average
			FROM students s
			JOIN grades g ON g.student_id = s.id
			GROUP BY s.class_name, g.quarter
		)
		SELECT 
			class_name,
			ROUND(AVG(quarter_average), 2) as average
		FROM class_quarter_averages
		GROUP BY class_name
		ORDER BY class_name
	`

	var performances []ClassPerformance
	err := database.DB.Select(&performances, query)
	if err != nil {
		log.Printf("Ошибка при получении успеваемости по классам: %v", err)
		http.Error(w, "Ошибка при получении успеваемости по классам: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Если нет данных, возвращаем пустой массив
	if performances == nil {
		performances = []ClassPerformance{}
	}

	json.NewEncoder(w).Encode(performances)
}
