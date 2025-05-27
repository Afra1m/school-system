package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"school-system/backend/database"
	"school-system/backend/models"
)

func GetSubjects(w http.ResponseWriter, r *http.Request) {
	log.Printf("Получен запрос на получение списка предметов")
	row, err := database.DB.Query(`SELECT DISTINCT id, name FROM subjects ORDER BY name`)
	if err != nil {
		log.Printf("Ошибка при получении данных предметов: %v", err)
		http.Error(w, "Ошибка при получении данных", http.StatusInternalServerError)
		return
	}
	defer row.Close()

	var subjects []models.Subject
	for row.Next() {
		var s models.Subject
		err := row.Scan(&s.ID, &s.Name)
		if err != nil {
			log.Printf("Ошибка при сканировании данных предмета: %v", err)
			http.Error(w, "Ошибка при получении данных", http.StatusInternalServerError)
			return
		}
		subjects = append(subjects, s)
	}

	log.Printf("Успешно получено %d предметов", len(subjects))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(subjects)
}

func CreateSubject(w http.ResponseWriter, r *http.Request) {
	log.Printf("Получен запрос на создание нового предмета")
	var subject models.Subject
	if err := json.NewDecoder(r.Body).Decode(&subject); err != nil {
		log.Printf("Ошибка при чтении данных предмета: %v", err)
		http.Error(w, "Невозможно прочитать данные", http.StatusBadRequest)
		return
	}

	_, err := database.DB.Exec(`INSERT INTO subjects (name, teacher_id) VALUES ($1, $2)`, subject.Name, subject.TeacherID)
	if err != nil {
		log.Printf("Ошибка при добавлении предмета в базу данных: %v", err)
		http.Error(w, "Ошибка при добавлении предмета", http.StatusInternalServerError)
		return
	}

	log.Printf("Успешно создан новый предмет: %s", subject.Name)
	w.WriteHeader(http.StatusCreated)
}

func UpdateSubject(w http.ResponseWriter, r *http.Request) {
	log.Printf("Получен запрос на обновление предмета")
	var s models.Subject
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		log.Printf("Ошибка при чтении данных предмета: %v", err)
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}

	_, err := database.DB.Exec(`UPDATE subjects SET name=$1 WHERE id=$2`, s.Name, s.ID)
	if err != nil {
		log.Printf("Ошибка при обновлении предмета с ID %d: %v", s.ID, err)
		http.Error(w, "Ошибка при обновлении", http.StatusInternalServerError)
		return
	}

	log.Printf("Успешно обновлен предмет с ID: %d", s.ID)
	w.WriteHeader(http.StatusOK)
}

func DeleteSubject(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	log.Printf("Получен запрос на удаление предмета с ID: %s", id)

	_, err := database.DB.Exec(`DELETE FROM subjects WHERE id=$1`, id)
	if err != nil {
		log.Printf("Ошибка при удалении предмета с ID %s: %v", id, err)
		http.Error(w, "Ошибка при удалении", http.StatusInternalServerError)
		return
	}

	log.Printf("Успешно удален предмет с ID: %s", id)
	w.WriteHeader(http.StatusOK)
}
