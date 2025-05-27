package middleware

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret []byte // теперь будет задаваться извне через SetJWTSecret

type contextKey string

const (
	ContextUserID contextKey = "userID"
	ContextRole   contextKey = "role"
)

// Функция для установки секрета JWT
func SetJWTSecret(secret []byte) {
	log.Printf("Установка секретного ключа JWT")
	jwtSecret = secret
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Проверка аутентификации для запроса: %s %s", r.Method, r.URL.Path)

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			log.Printf("Отсутствует заголовок Authorization")
			http.Error(w, "Нет токена", http.StatusUnauthorized)
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		log.Printf("Получен токен: %s...", tokenStr[:10])

		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil {
			log.Printf("Ошибка при проверке токена: %v", err)
			http.Error(w, "Неверный токен", http.StatusUnauthorized)
			return
		}

		if !token.Valid {
			log.Printf("Токен недействителен")
			http.Error(w, "Неверный токен", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			log.Printf("Неверный формат claims токена")
			http.Error(w, "Неверные claims", http.StatusUnauthorized)
			return
		}

		userID := claims["user_id"]
		role := claims["role"]
		log.Printf("Успешная аутентификация: user_id=%v, role=%v", userID, role)

		// Добавим userID и role в context
		ctx := context.WithValue(r.Context(), ContextUserID, userID)
		ctx = context.WithValue(ctx, ContextRole, role)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
