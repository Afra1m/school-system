package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"school-system/backend/database"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// Получаем секрет JWT из окружения или используем фиксированный для тестирования
func getJWTSecret() []byte {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Printf("ВНИМАНИЕ: Используется тестовый секретный ключ JWT")
		jwtSecret = "test-secret-key-123" // Фиксированный ключ для тестирования
	}
	return []byte(jwtSecret)
}

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Login(w http.ResponseWriter, r *http.Request) {
	log.Printf("Получен запрос на вход в систему")
	var creds Credentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		log.Printf("Ошибка при чтении данных входа: %v", err)
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}

	user, err := database.GetUserByUsername(creds.Username)
	if err != nil {
		log.Printf("Пользователь не найден: %s", creds.Username)
		http.Error(w, "Пользователь не найден", http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password)); err != nil {
		log.Printf("Неверный пароль для пользователя: %s", creds.Username)
		http.Error(w, "Неверный пароль", http.StatusUnauthorized)
		return
	}

	// Генерация токена
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString(getJWTSecret())
	if err != nil {
		log.Printf("Ошибка при генерации токена для пользователя %s: %v", creds.Username, err)
		http.Error(w, "Ошибка генерации токена", http.StatusInternalServerError)
		return
	}

	log.Printf("Успешный вход пользователя: %s (роль: %s)", creds.Username, user.Role)
	json.NewEncoder(w).Encode(map[string]string{
		"token": tokenString,
	})
}

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"` // student / teacher / deputy
}

func Register(w http.ResponseWriter, r *http.Request) {
	log.Printf("Получен запрос на регистрацию нового пользователя")
	var req RegisterRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Printf("Ошибка при чтении данных регистрации: %v", err)
		http.Error(w, "Некорректные данные", http.StatusBadRequest)
		return
	}

	// Проверяем роль
	if req.Role != "student" && req.Role != "teacher" && req.Role != "deputy" {
		log.Printf("Попытка регистрации с недопустимой ролью: %s", req.Role)
		http.Error(w, "Недопустимая роль. Допустимые значения: student, teacher, deputy", http.StatusBadRequest)
		return
	}

	// Проверим, есть ли уже пользователь с таким именем
	existingUser, _ := database.GetUserByUsername(req.Username)
	if existingUser != nil && existingUser.Username != "" {
		log.Printf("Попытка регистрации существующего пользователя: %s", req.Username)
		http.Error(w, "Пользователь уже существует", http.StatusConflict)
		return
	}

	// Хешируем пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Ошибка при хешировании пароля для пользователя %s: %v", req.Username, err)
		http.Error(w, "Ошибка при хешировании пароля", http.StatusInternalServerError)
		return
	}

	// Сохраняем пользователя в базу
	err = database.CreateUser(req.Username, string(hashedPassword), req.Role)
	if err != nil {
		log.Printf("Ошибка при создании пользователя %s: %v", req.Username, err)
		http.Error(w, "Ошибка при создании пользователя", http.StatusInternalServerError)
		return
	}

	log.Printf("Успешная регистрация нового пользователя: %s (роль: %s)", req.Username, req.Role)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Пользователь зарегистрирован",
	})
}

func VerifyToken(w http.ResponseWriter, r *http.Request) {
	log.Printf("Получен запрос на проверку токена")

	// Получаем токен из заголовка Authorization
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		log.Printf("Отсутствует заголовок Authorization")
		http.Error(w, "Требуется авторизация", http.StatusUnauthorized)
		return
	}

	// Убираем префикс "Bearer " из токена
	tokenString := authHeader
	if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
		tokenString = tokenString[7:]
	}

	// Парсим и проверяем токен
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("неожиданный метод подписи: %v", token.Header["alg"])
		}
		return getJWTSecret(), nil
	})

	if err != nil {
		log.Printf("Ошибка при проверке токена: %v", err)
		http.Error(w, "Недействительный токен", http.StatusUnauthorized)
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Проверяем срок действия токена
		if exp, ok := claims["exp"].(float64); ok {
			if time.Now().Unix() > int64(exp) {
				log.Printf("Токен истек")
				http.Error(w, "Токен истек", http.StatusUnauthorized)
				return
			}
		}

		// Токен действителен
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"valid":   true,
			"user_id": claims["user_id"],
			"role":    claims["role"],
		})
		return
	}

	log.Printf("Недействительный токен")
	http.Error(w, "Недействительный токен", http.StatusUnauthorized)
}
