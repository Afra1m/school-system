package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"

	"school-system/backend/models" // Используй свой путь
)

// newRouter — локальный роутер для тестов, с нужными маршрутами.
func newRouter() *mux.Router {
	r := mux.NewRouter()

	// Маршруты для студентов (которые нужны в тестах)
	r.HandleFunc("/students", GetStudents).Methods("GET")
	r.HandleFunc("/students", CreateStudent).Methods("POST")
	r.HandleFunc("/students/{id}", UpdateStudent).Methods("PUT")
	r.HandleFunc("/students/{id}", DeleteStudent).Methods("DELETE")

	return r
}

func TestCreateStudent(t *testing.T) {
	router := newRouter()
	newStudent := models.Student{
		FullName:  "Иван Иванов",
		ClassName: "9А",
	}

	body, _ := json.Marshal(newStudent)
	req, _ := http.NewRequest("POST", "/students", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusCreated {
		t.Errorf("Ожидался статус 201, получен %d", resp.Code)
	}
}

func TestGetStudents(t *testing.T) {
	router := newRouter()

	req, _ := http.NewRequest("GET", "/students", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("Ожидался статус 200, получен %d", resp.Code)
	}

	var students []models.Student
	err := json.Unmarshal(resp.Body.Bytes(), &students)
	if err != nil {
		t.Errorf("Ошибка при разборе ответа: %v", err)
	}
}

func TestUpdateStudent(t *testing.T) {
	router := newRouter()

	updatedStudent := models.Student{
		FullName:  "Пётр Петров",
		ClassName: "10Б",
	}

	body, _ := json.Marshal(updatedStudent)
	req, _ := http.NewRequest("PUT", "/students/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("Ожидался статус 200, получен %d", resp.Code)
	}
}

func TestDeleteStudent(t *testing.T) {
	router := newRouter()

	req, _ := http.NewRequest("DELETE", "/students/1", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("Ожидался статус 200, получен %d", resp.Code)
	}
}
