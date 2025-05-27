package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"school-system/backend/database"
	"school-system/backend/models"

	"github.com/gorilla/mux"
)

func GetStudents(w http.ResponseWriter, r *http.Request) {
	log.Printf("Получен запрос на получение списка студентов")
	rows, err := database.DB.Query("SELECT id, full_name, class_name FROM students")
	if err != nil {
		log.Printf("Ошибка при получении данных студентов: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var students []models.Student
	for rows.Next() {
		var student models.Student
		if err := rows.Scan(&student.ID, &student.FullName, &student.ClassName); err != nil {
			log.Printf("Ошибка при сканировании данных студента: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		students = append(students, student)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Ошибка при обработке результатов запроса: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Успешно получено %d студентов", len(students))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(students)
}

func CreateStudent(w http.ResponseWriter, r *http.Request) {
	log.Printf("Получен запрос на создание нового студента")
	var student models.Student
	if err := json.NewDecoder(r.Body).Decode(&student); err != nil {
		log.Printf("Ошибка при чтении данных студента: %v", err)
		http.Error(w, "Невозможно прочитать данные", http.StatusBadRequest)
		return
	}

	_, err := database.DB.Exec(`INSERT INTO students (full_name, class_name) VALUES ($1, $2)`,
		student.FullName, student.ClassName)
	if err != nil {
		log.Printf("Ошибка при добавлении студента в базу данных: %v", err)
		http.Error(w, "Ошибка при добавлении студента", http.StatusInternalServerError)
		return
	}

	log.Printf("Успешно создан новый студент: %s", student.FullName)
	w.WriteHeader(http.StatusCreated)
}

func DeleteStudent(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	log.Printf("Получен запрос на удаление студента с ID: %s", id)

	_, err := database.DB.Exec(`DELETE FROM students WHERE id = $1`, id)
	if err != nil {
		log.Printf("Ошибка при удалении студента с ID %s: %v", id, err)
		http.Error(w, "Ошибка при удалении студента", http.StatusInternalServerError)
		return
	}

	log.Printf("Успешно удален студент с ID: %s", id)
	w.WriteHeader(http.StatusNoContent)
}

func UpdateStudent(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	log.Printf("Получен запрос на обновление студента с ID: %s", id)

	var student models.Student
	if err := json.NewDecoder(r.Body).Decode(&student); err != nil {
		log.Printf("Ошибка при чтении данных студента: %v", err)
		http.Error(w, "Невозможно прочитать данные", http.StatusBadRequest)
		return
	}

	_, err := database.DB.Exec(`UPDATE students SET full_name = $1, class_name = $2 WHERE id = $3`,
		student.FullName, student.ClassName, id)
	if err != nil {
		log.Printf("Ошибка при обновлении студента с ID %s: %v", id, err)
		http.Error(w, "Ошибка при обновлении студента", http.StatusInternalServerError)
		return
	}

	log.Printf("Успешно обновлен студент с ID: %s", id)
	w.WriteHeader(http.StatusOK)
}

// GetFailingStudents возвращает список неуспевающих учеников с их средними оценками по предметам
func GetFailingStudents(w http.ResponseWriter, r *http.Request) {
	log.Printf("Получен запрос на получение списка отстающих учеников")

	students, err := database.GetFailingStudents()
	if err != nil {
		log.Printf("Ошибка при получении списка неуспевающих учеников: %v", err)
		http.Error(w, fmt.Sprintf("Ошибка при получении списка неуспевающих учеников: %v", err), http.StatusInternalServerError)
		return
	}

	log.Printf("Найдено отстающих учеников: %d", len(students))
	for _, student := range students {
		log.Printf("Отстающий ученик: %s (Класс: %s), предметов: %d",
			student.FullName,
			student.ClassName,
			len(student.SubjectAverages))
		for _, subject := range student.SubjectAverages {
			log.Printf("  - %s: %.2f", subject.SubjectName, subject.Average)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(students)
}

// GetAverageGradesByClass возвращает средние оценки по предметам для каждого класса
func GetAverageGradesByClass(w http.ResponseWriter, r *http.Request) {
	averages, err := database.GetAverageGradesByClass()
	if err != nil {
		http.Error(w, "Ошибка при получении средних оценок", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(averages)
}

// GetTopAndWorstClasses возвращает классы с самой высокой и низкой успеваемостью
func GetTopAndWorstClasses(w http.ResponseWriter, r *http.Request) {
	topClass, worstClass, err := database.GetTopAndWorstClasses()
	if err != nil {
		http.Error(w, "Ошибка при получении информации о классах", http.StatusInternalServerError)
		return
	}

	response := struct {
		TopClass   string `json:"top_class"`
		WorstClass string `json:"worst_class"`
	}{
		TopClass:   topClass,
		WorstClass: worstClass,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
