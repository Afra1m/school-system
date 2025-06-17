package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"school-system/backend/database"
	"school-system/backend/middleware"
	"school-system/backend/models"
	"strconv"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

func GetTeachers(w http.ResponseWriter, r *http.Request) {
	log.Printf("Получен запрос на получение списка учителей")
	row, err := database.DB.Query(`SELECT id, full_name, room_number, user_id FROM teachers`)
	if err != nil {
		log.Printf("Ошибка при получении данных учителей: %v", err)
		http.Error(w, "Ошибка при получении данных", http.StatusInternalServerError)
		return
	}
	defer row.Close()

	var teachers []models.Teacher

	for row.Next() {
		var t models.Teacher

		err := row.Scan(&t.ID, &t.FullName, &t.RoomNumber, &t.UserID)
		if err != nil {
			log.Printf("Ошибка при сканировании данных учителя: %v", err)
			http.Error(w, "Ошибка при получении данных: "+err.Error(), http.StatusInternalServerError)
			return
		}
		teachers = append(teachers, t)
	}

	log.Printf("Успешно получено %d учителей", len(teachers))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(teachers)
}

func CreateTeacher(w http.ResponseWriter, r *http.Request) {
	log.Printf("Получен запрос на создание нового учителя")
	var teacher models.Teacher
	if err := json.NewDecoder(r.Body).Decode(&teacher); err != nil {
		log.Printf("Ошибка при чтении данных учителя: %v", err)
		http.Error(w, "Невозможно прочитать данные", http.StatusBadRequest)
		return
	}

	// Создаем пользователя для учителя
	username := teacher.FullName // Используем ФИО как логин
	password := "password123"    // Временный пароль, который учитель должен будет сменить
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Ошибка при хешировании пароля: %v", err)
		http.Error(w, "Ошибка при создании пользователя", http.StatusInternalServerError)
		return
	}

	// Создаем пользователя
	var userID int
	err = database.DB.QueryRow(
		`INSERT INTO users (username, password, role) VALUES ($1, $2, $3) RETURNING id`,
		username, string(hashedPassword), "teacher",
	).Scan(&userID)
	if err != nil {
		log.Printf("Ошибка при создании пользователя: %v", err)
		http.Error(w, "Ошибка при создании пользователя", http.StatusInternalServerError)
		return
	}

	// Создаем учителя
	_, err = database.DB.Exec(`INSERT INTO teachers (full_name, room_number, user_id) VALUES ($1, $2, $3)`,
		teacher.FullName, teacher.RoomNumber, userID)
	if err != nil {
		log.Printf("Ошибка при добавлении учителя в базу данных: %v", err)
		http.Error(w, "Ошибка при добавлении учителя", http.StatusInternalServerError)
		return
	}

	log.Printf("Успешно создан новый учитель: %s", teacher.FullName)
	w.WriteHeader(http.StatusCreated)
}

func DeleteTeacher(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	log.Printf("Получен запрос на удаление учителя с ID: %s", id)

	_, err := database.DB.Exec(`DELETE FROM teachers WHERE id = $1`, id)
	if err != nil {
		log.Printf("Ошибка при удалении учителя с ID %s: %v", id, err)
		http.Error(w, "Ошибка при удалении учителя", http.StatusInternalServerError)
		return
	}

	log.Printf("Успешно удален учитель с ID: %s", id)
	w.WriteHeader(http.StatusNoContent)
}

func UpdateTeacher(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	log.Printf("Получен запрос на обновление учителя с ID: %s", id)

	var teacher models.Teacher
	if err := json.NewDecoder(r.Body).Decode(&teacher); err != nil {
		log.Printf("Ошибка при чтении данных учителя: %v", err)
		http.Error(w, "Невозможно прочитать данные", http.StatusBadRequest)
		return
	}

	_, err := database.DB.Exec(`UPDATE teachers SET full_name = $1, room_number = $2 WHERE id = $3`,
		teacher.FullName, teacher.RoomNumber, id)
	if err != nil {
		log.Printf("Ошибка при обновлении учителя с ID %s: %v", id, err)
		http.Error(w, "Ошибка при обновлении данных учителя", http.StatusInternalServerError)
		return
	}

	log.Printf("Успешно обновлен учитель с ID: %s", id)
	w.WriteHeader(http.StatusOK)
}

// GetMyStudents возвращает список учеников для конкретного учителя
func GetMyStudents(w http.ResponseWriter, r *http.Request) {
	// Получаем ID учителя из JWT токена
	idVal := r.Context().Value(middleware.ContextUserID)
	log.Printf("Получено значение из контекста: %v (тип: %T)", idVal, idVal)

	var teacherID int
	switch v := idVal.(type) {
	case float64:
		teacherID = int(v)
		log.Printf("Преобразовано из float64: %d", teacherID)
	case int:
		teacherID = v
		log.Printf("Получено как int: %d", teacherID)
	case string:
		// если вдруг строка, попробуем преобразовать
		if id, err := strconv.Atoi(v); err == nil {
			teacherID = id
			log.Printf("Преобразовано из string: %d", teacherID)
		} else {
			log.Printf("Ошибка преобразования строки в число: %v", err)
			http.Error(w, "Некорректный user_id в токене", http.StatusUnauthorized)
			return
		}
	default:
		log.Printf("Неизвестный тип данных: %T", v)
		http.Error(w, "Не удалось определить user_id", http.StatusUnauthorized)
		return
	}

	log.Printf("Используем teacherID=%d для поиска предметов", teacherID)

	// Получаем предметы, которые ведет учитель
	subjects, err := database.GetSubjectsByTeacher(teacherID)
	if err != nil {
		log.Printf("Ошибка при получении предметов: %v", err)
		http.Error(w, "Ошибка при получении предметов учителя", http.StatusInternalServerError)
		return
	}

	log.Printf("Найдено предметов: %d", len(subjects))

	// Получаем всех учеников, которые изучают эти предметы
	var students []models.Student
	for _, subject := range subjects {
		log.Printf("Ищем учеников для предмета id=%d, name=%s", subject.ID, subject.Name)
		subjectStudents, err := database.GetStudentsBySubject(subject.ID)
		if err != nil {
			log.Printf("Ошибка при получении учеников для предмета %d: %v", subject.ID, err)
			continue
		}
		log.Printf("Найдено учеников для предмета %d: %d", subject.ID, len(subjectStudents))
		students = append(students, subjectStudents...)
	}

	// Удаляем дубликаты
	uniqueStudents := make(map[int]models.Student)
	for _, student := range students {
		uniqueStudents[student.ID] = student
	}

	var result []models.Student
	for _, student := range uniqueStudents {
		result = append(result, student)
	}

	log.Printf("Итоговое количество уникальных учеников: %d", len(result))
	json.NewEncoder(w).Encode(result)
}

// GetMyStudentsGrades возвращает оценки учеников конкретного учителя
func GetMyStudentsGrades(w http.ResponseWriter, r *http.Request) {
	// Получаем ID учителя из JWT токена
	idVal := r.Context().Value(middleware.ContextUserID)
	var teacherID int
	switch v := idVal.(type) {
	case float64:
		teacherID = int(v)
	case int:
		teacherID = v
	case string:
		if id, err := strconv.Atoi(v); err == nil {
			teacherID = id
		} else {
			http.Error(w, "Некорректный user_id в токене", http.StatusUnauthorized)
			return
		}
	default:
		http.Error(w, "Не удалось определить user_id", http.StatusUnauthorized)
		return
	}

	// Получаем оценки учеников
	studentGrades, err := database.GetGradesByTeacherAndStudents(teacherID)
	if err != nil {
		http.Error(w, "Ошибка при получении оценок учеников", http.StatusInternalServerError)
		return
	}

	// Формируем ответ
	type StudentGrades struct {
		StudentID   int            `json:"student_id"`
		StudentName string         `json:"student_name"`
		Grades      []models.Grade `json:"grades"`
	}

	var response []StudentGrades
	for studentID, grades := range studentGrades {
		// Получаем имя ученика
		var student models.Student
		err := database.DB.Get(&student, "SELECT full_name FROM students WHERE id = $1", studentID)
		if err != nil {
			continue
		}

		response = append(response, StudentGrades{
			StudentID:   studentID,
			StudentName: student.FullName,
			Grades:      grades,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
