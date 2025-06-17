package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"

	"school-system/backend/database"
	"school-system/backend/handlers"
	"school-system/backend/middleware"
)

func main() {
	log.Printf("Запуск сервера...")

	// Загружаем .env
	err := godotenv.Load()
	if err != nil {
		log.Printf("Предупреждение: .env файл не найден, используются переменные окружения")
	}

	// Получаем секрет JWT из окружения или используем фиксированный для тестирования
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Printf("Предупреждение: используется тестовый секретный ключ JWT")
		jwtSecret = "test-secret-key-123" // Фиксированный ключ для тестирования
	}

	// Передаем секрет в middleware
	middleware.SetJWTSecret([]byte(jwtSecret))
	log.Printf("Секретный ключ JWT установлен")

	// Инициализируем базу данных
	log.Printf("Инициализация подключения к базе данных...")
	database.InitDB()

	// Настройка роутеров
	log.Printf("Настройка маршрутов...")
	r := mux.NewRouter()

	// Настройка CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization", "Accept", "Origin", "X-Requested-With"},
		ExposedHeaders:   []string{"Content-Length"},
		AllowCredentials: true,
		Debug:            true,
	})

	// ====== GET Requests ======
	log.Printf("Регистрация GET маршрутов...")
	r.HandleFunc("/students", handlers.GetStudents).Methods("GET")
	r.HandleFunc("/teachers", handlers.GetTeachers).Methods("GET")
	r.HandleFunc("/subjects", handlers.GetSubjects).Methods("GET")

	// Получение оценок — только для авторизованных пользователей
	r.Handle("/grades", middleware.AuthMiddleware(http.HandlerFunc(handlers.GetGrades))).Methods("GET")

	// Получение оценок конкретного студента
	r.Handle("/grades/student/{id}", middleware.AuthMiddleware(http.HandlerFunc(handlers.GetStudentGrades))).Methods("GET")

	// Новые маршруты для завуча
	r.Handle("/students/failing", middleware.AuthMiddleware(
		http.HandlerFunc(
			middleware.RequireRole("deputy")(handlers.GetFailingStudents),
		),
	)).Methods("GET")

	r.Handle("/grades/average-by-class", middleware.AuthMiddleware(
		http.HandlerFunc(
			middleware.RequireRole("deputy")(handlers.GetAverageGradesByClass),
		),
	)).Methods("GET")

	r.Handle("/stats/top-worst-classes", middleware.AuthMiddleware(
		http.HandlerFunc(
			middleware.RequireRole("deputy")(handlers.GetTopAndWorstClasses),
		),
	)).Methods("GET")

	// Маршрут для учителя
	r.Handle("/teacher/my-students", middleware.AuthMiddleware(
		http.HandlerFunc(
			middleware.RequireRole("teacher")(handlers.GetMyStudents),
		),
	)).Methods("GET")

	r.Handle("/teacher/my-students/grades", middleware.AuthMiddleware(
		http.HandlerFunc(
			middleware.RequireRole("teacher")(handlers.GetMyStudentsGrades),
		),
	)).Methods("GET")

	// Новые маршруты для статистики
	r.Handle("/stats/students-count", middleware.AuthMiddleware(
		http.HandlerFunc(handlers.GetStudentsCount),
	)).Methods("GET")

	r.Handle("/stats/teachers-count", middleware.AuthMiddleware(
		http.HandlerFunc(handlers.GetTeachersCount),
	)).Methods("GET")

	r.Handle("/stats/average-grade", middleware.AuthMiddleware(
		http.HandlerFunc(handlers.GetAverageGrade),
	)).Methods("GET")

	r.Handle("/stats/class-performance", middleware.AuthMiddleware(
		http.HandlerFunc(handlers.GetClassPerformance),
	)).Methods("GET")

	r.Handle("/stats/average-grades", middleware.AuthMiddleware(
		http.HandlerFunc(
			middleware.RequireRole("deputy")(handlers.GetAverageGradesByClass),
		),
	)).Methods("GET")

	r.Handle("/stats/failing-students", middleware.AuthMiddleware(
		http.HandlerFunc(
			middleware.RequireRole("deputy")(handlers.GetFailingStudents),
		),
	)).Methods("GET")

	r.Handle("/stats/top-worst-classes", middleware.AuthMiddleware(
		http.HandlerFunc(
			middleware.RequireRole("deputy")(handlers.GetTopAndWorstClasses),
		),
	)).Methods("GET")

	// ====== POST Requests ======
	log.Printf("Регистрация POST маршрутов...")
	r.HandleFunc("/students", handlers.CreateStudent).Methods("POST")
	r.HandleFunc("/subjects", handlers.CreateSubject).Methods("POST")
	r.Handle("/teachers", middleware.AuthMiddleware(
		http.HandlerFunc(
			middleware.RequireRole("deputy")(handlers.CreateTeacher),
		),
	)).Methods("POST")

	// Создание оценки — только для ролей "deputy" и "teacher", с авторизацией
	r.Handle("/grades", middleware.AuthMiddleware(
		http.HandlerFunc(
			middleware.RequireRole("deputy", "teacher")(handlers.CreateGrade),
		),
	)).Methods("POST")

	// ====== PUT Requests ======
	log.Printf("Регистрация PUT маршрутов...")
	r.HandleFunc("/students/{id}", handlers.UpdateStudent).Methods("PUT")
	r.HandleFunc("/teachers/{id}", handlers.UpdateTeacher).Methods("PUT")
	r.HandleFunc("/subjects/{id}", handlers.UpdateSubject).Methods("PUT")
	r.Handle("/grades/{id}", middleware.AuthMiddleware(
		http.HandlerFunc(
			middleware.RequireRole("deputy", "teacher")(handlers.UpdateGrade),
		),
	)).Methods("PUT")

	// ====== DELETE Requests ======
	log.Printf("Регистрация DELETE маршрутов...")
	r.Handle("/students/{id}", middleware.AuthMiddleware(
		http.HandlerFunc(
			middleware.RequireRole("deputy")(handlers.DeleteStudent),
		),
	)).Methods("DELETE")
	r.Handle("/teachers/{id}", middleware.AuthMiddleware(
		http.HandlerFunc(
			middleware.RequireRole("deputy")(handlers.DeleteTeacher),
		),
	)).Methods("DELETE")
	r.HandleFunc("/subjects/{id}", handlers.DeleteSubject).Methods("DELETE")

	// Удаление оценки — только для роли "deputy", с авторизацией
	r.Handle("/grades/{id}", middleware.AuthMiddleware(
		http.HandlerFunc(
			middleware.RequireRole("deputy")(handlers.DeleteGrade),
		),
	)).Methods("DELETE")

	// ====== Аутентификация ======
	log.Printf("Регистрация маршрутов аутентификации...")
	r.HandleFunc("/login", handlers.Login).Methods("POST")
	r.HandleFunc("/register", handlers.Register).Methods("POST")
	r.HandleFunc("/verify-token", handlers.VerifyToken).Methods("GET")

	// Запуск сервера с CORS middleware
	log.Printf("Сервер запущен на порту 8000")
	handler := c.Handler(r)
	log.Fatal(http.ListenAndServe(":8000", handler))
}
