package database

import (
	"fmt"
	"log"

	"school-system/backend/models"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var DB *sqlx.DB

func InitDB() {
	var err error
	connStr := "host=localhost port=5433 user=postgres password=12345 dbname=school_system sslmode=disable"
	DB, err = sqlx.Connect("postgres", connStr)
	if err != nil {
		log.Fatal("Ошибка подключения к базе данных: ", err)
	}
	if err := DB.Ping(); err != nil {
		log.Fatal("Не удалось подключиться к базе данных: ", err)
	}
	fmt.Println("Подключение к базе данных успешно!")
}

func GetUserByUsername(username string) (*models.User, error) {
	user := &models.User{}
	row := DB.QueryRowx("SELECT id, username, password, role FROM users WHERE username = $1", username)
	err := row.StructScan(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// Создать нового пользователя
func CreateUser(username, hashedPassword, role string) error {
	_, err := DB.Exec("INSERT INTO users (username, password, role) VALUES ($1, $2, $3)",
		username, hashedPassword, role)
	return err
}

// GetGradesByTeacher возвращает все оценки по предметам конкретного учителя
func GetGradesByTeacher(teacherID int) ([]models.Grade, error) {
	var grades []models.Grade
	query := `
		SELECT g.* FROM grades g
		JOIN subjects s ON g.subject_id = s.id
		WHERE s.teacher_id = $1
	`
	err := DB.Select(&grades, query, teacherID)
	return grades, err
}

// GetTeacherByUserID возвращает учителя по его user_id
func GetTeacherByUserID(userID int) (*models.Teacher, error) {
	log.Printf("Ищем учителя с user_id=%d", userID)
	var teacher models.Teacher
	query := "SELECT id, full_name, room_number, user_id FROM teachers WHERE user_id = $1"
	log.Printf("Выполняем запрос: %s с параметром user_id=%d", query, userID)

	err := DB.Get(&teacher, query, userID)
	if err != nil {
		log.Printf("Ошибка при поиске учителя: %v", err)
		return nil, err
	}
	log.Printf("Учитель найден: id=%d, full_name=%s", teacher.ID, teacher.FullName)
	return &teacher, nil
}

// GetSubjectsByTeacher возвращает все предметы, которые ведет учитель
func GetSubjectsByTeacher(userID int) ([]models.Subject, error) {
	log.Printf("Получаем предметы для user_id=%d", userID)

	// Сначала получаем ID учителя из таблицы teachers
	teacher, err := GetTeacherByUserID(userID)
	if err != nil {
		log.Printf("Ошибка при получении учителя: %v", err)
		return nil, err
	}
	log.Printf("Найден учитель: id=%d, full_name=%s", teacher.ID, teacher.FullName)

	var subjects []models.Subject
	query := `SELECT * FROM subjects WHERE teacher_id = $1`
	err = DB.Select(&subjects, query, teacher.ID)
	if err != nil {
		log.Printf("Ошибка при получении предметов: %v", err)
		return nil, err
	}
	log.Printf("Найдено предметов: %d", len(subjects))
	return subjects, err
}

// GetStudentsBySubject возвращает всех учеников, изучающих конкретный предмет
func GetStudentsBySubject(subjectID int) ([]models.Student, error) {
	log.Printf("Ищем учеников для предмета с id=%d", subjectID)
	var students []models.Student
	query := `
		SELECT DISTINCT s.* FROM students s
		JOIN grades g ON s.id = g.student_id
		WHERE g.subject_id = $1
	`
	log.Printf("Выполняем запрос: %s с параметром subject_id=%d", query, subjectID)

	err := DB.Select(&students, query, subjectID)
	if err != nil {
		log.Printf("Ошибка при поиске учеников: %v", err)
		return nil, err
	}
	log.Printf("Найдено учеников: %d", len(students))
	return students, err
}

// GetFailingStudents возвращает список неуспевающих учеников с их средними оценками по предметам
func GetFailingStudents() ([]struct {
	models.Student
	SubjectAverages []struct {
		SubjectName string  `db:"subject_name" json:"subject_name"`
		Quarter     int     `db:"quarter" json:"quarter"`
		Average     float64 `db:"average" json:"average"`
	} `json:"subject_averages"`
}, error) {
	query := `
		WITH student_subject_averages AS (
			SELECT 
				s.id,
				s.full_name,
				s.class_name,
				s.user_id,
				sub.name as subject_name,
				g.quarter,
				AVG(g.grade) as average
			FROM students s
			JOIN grades g ON s.id = g.student_id
			JOIN subjects sub ON g.subject_id = sub.id
			GROUP BY s.id, s.full_name, s.class_name, s.user_id, sub.name, g.quarter
			HAVING AVG(g.grade) < 3
		)
		SELECT 
			id,
			full_name,
			class_name,
			user_id,
			subject_name,
			quarter,
			average
		FROM student_subject_averages
		ORDER BY full_name, subject_name, quarter
	`
	log.Printf("Выполняем запрос для получения отстающих студентов: %s", query)

	rows, err := DB.Queryx(query)
	if err != nil {
		log.Printf("Ошибка при выполнении запроса: %v", err)
		return nil, err
	}
	defer rows.Close()

	// Мапа для группировки данных по студентам
	studentMap := make(map[int]struct {
		models.Student
		SubjectAverages []struct {
			SubjectName string  `db:"subject_name" json:"subject_name"`
			Quarter     int     `db:"quarter" json:"quarter"`
			Average     float64 `db:"average" json:"average"`
		} `json:"subject_averages"`
	})

	for rows.Next() {
		var student models.Student
		var subjectName string
		var quarter int
		var average float64

		if err := rows.Scan(&student.ID, &student.FullName, &student.ClassName, &student.UserID, &subjectName, &quarter, &average); err != nil {
			log.Printf("Ошибка при сканировании строки: %v", err)
			return nil, err
		}

		log.Printf("Найден отстающий студент: ID=%d, ФИО=%s, Класс=%s, Предмет=%s, Четверть=%d, Средний балл=%.2f",
			student.ID, student.FullName, student.ClassName, subjectName, quarter, average)

		// Если студент еще не в мапе, добавляем его
		if _, exists := studentMap[student.ID]; !exists {
			studentMap[student.ID] = struct {
				models.Student
				SubjectAverages []struct {
					SubjectName string  `db:"subject_name" json:"subject_name"`
					Quarter     int     `db:"quarter" json:"quarter"`
					Average     float64 `db:"average" json:"average"`
				} `json:"subject_averages"`
			}{
				Student: student,
				SubjectAverages: make([]struct {
					SubjectName string  `db:"subject_name" json:"subject_name"`
					Quarter     int     `db:"quarter" json:"quarter"`
					Average     float64 `db:"average" json:"average"`
				}, 0),
			}
		}

		// Добавляем информацию о предмете
		studentData := studentMap[student.ID]
		studentData.SubjectAverages = append(studentData.SubjectAverages, struct {
			SubjectName string  `db:"subject_name" json:"subject_name"`
			Quarter     int     `db:"quarter" json:"quarter"`
			Average     float64 `db:"average" json:"average"`
		}{
			SubjectName: subjectName,
			Quarter:     quarter,
			Average:     average,
		})
		studentMap[student.ID] = studentData
	}

	if err := rows.Err(); err != nil {
		log.Printf("Ошибка при итерации по строкам: %v", err)
		return nil, err
	}

	// Преобразуем мапу в слайс
	var result []struct {
		models.Student
		SubjectAverages []struct {
			SubjectName string  `db:"subject_name" json:"subject_name"`
			Quarter     int     `db:"quarter" json:"quarter"`
			Average     float64 `db:"average" json:"average"`
		} `json:"subject_averages"`
	}

	for _, studentData := range studentMap {
		result = append(result, studentData)
	}

	log.Printf("Итоговое количество отстающих учеников: %d", len(result))
	return result, nil
}

// GetAverageGradesByClass возвращает средние оценки по предметам для каждого класса
func GetAverageGradesByClass() (map[string]map[string]float64, error) {
	result := make(map[string]map[string]float64)

	query := `
		WITH subject_quarter_averages AS (
			SELECT 
				s.class_name,
				sub.name as subject_name,
				g.quarter,
				AVG(g.grade) as quarter_average
			FROM students s
			JOIN grades g ON s.id = g.student_id
			JOIN subjects sub ON g.subject_id = sub.id
			GROUP BY s.class_name, sub.name, g.quarter
		)
		SELECT 
			class_name,
			subject_name,
			ROUND(AVG(quarter_average), 2) as average_grade
		FROM subject_quarter_averages
		GROUP BY class_name, subject_name
		ORDER BY class_name, subject_name
	`

	rows, err := DB.Queryx(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var className, subjectName string
		var averageGrade float64
		if err := rows.Scan(&className, &subjectName, &averageGrade); err != nil {
			return nil, err
		}

		if _, ok := result[className]; !ok {
			result[className] = make(map[string]float64)
		}
		result[className][subjectName] = averageGrade
	}

	return result, nil
}

// GetTopAndWorstClasses возвращает классы с самой высокой и низкой успеваемостью
func GetTopAndWorstClasses() (string, string, error) {
	query := `
		SELECT 
			s.class_name,
			AVG(g.grade) as class_average
		FROM students s
		JOIN grades g ON s.id = g.student_id
		GROUP BY s.class_name
		ORDER BY class_average DESC
	`

	rows, err := DB.Queryx(query)
	if err != nil {
		return "", "", err
	}
	defer rows.Close()

	var classes []struct {
		ClassName string
		Average   float64
	}

	for rows.Next() {
		var class struct {
			ClassName string
			Average   float64
		}
		if err := rows.Scan(&class.ClassName, &class.Average); err != nil {
			return "", "", err
		}
		classes = append(classes, class)
	}

	if len(classes) == 0 {
		return "", "", nil
	}

	return classes[0].ClassName, classes[len(classes)-1].ClassName, nil
}

// GetAllTeachers возвращает список всех учителей
func GetAllTeachers() ([]models.Teacher, error) {
	var teachers []models.Teacher
	query := `SELECT id, full_name, room_number, user_id FROM teachers`
	err := DB.Select(&teachers, query)
	return teachers, err
}

// GetGradesByTeacherAndStudents возвращает оценки учеников конкретного учителя
func GetGradesByTeacherAndStudents(userID int) (map[int][]models.Grade, error) {
	log.Printf("Получаем оценки для учителя с user_id=%d", userID)

	// Сначала получаем ID учителя из таблицы teachers
	teacher, err := GetTeacherByUserID(userID)
	if err != nil {
		log.Printf("Ошибка при получении учителя: %v", err)
		return nil, err
	}
	log.Printf("Найден учитель: id=%d, full_name=%s", teacher.ID, teacher.FullName)

	query := `
		SELECT g.*, s.full_name as student_name
		FROM grades g
		JOIN students s ON g.student_id = s.id
		JOIN subjects sub ON g.subject_id = sub.id
		WHERE sub.teacher_id = $1
		ORDER BY s.full_name, g.quarter
	`
	log.Printf("Выполняем запрос: %s с параметром teacher_id=%d", query, teacher.ID)

	rows, err := DB.Queryx(query, teacher.ID)
	if err != nil {
		log.Printf("Ошибка при выполнении запроса: %v", err)
		return nil, err
	}
	defer rows.Close()

	// Мапа для хранения оценок по ученикам
	studentGrades := make(map[int][]models.Grade)

	for rows.Next() {
		var grade models.Grade
		var studentName string
		if err := rows.Scan(&grade.ID, &grade.StudentID, &grade.SubjectID, &grade.Grade, &grade.Quarter, &studentName); err != nil {
			log.Printf("Ошибка при сканировании строки: %v", err)
			return nil, err
		}
		log.Printf("Найдена оценка: student_id=%d, subject_id=%d, grade=%d", grade.StudentID, grade.SubjectID, grade.Grade)
		studentGrades[grade.StudentID] = append(studentGrades[grade.StudentID], grade)
	}

	log.Printf("Найдено оценок для %d учеников", len(studentGrades))
	return studentGrades, nil
}
