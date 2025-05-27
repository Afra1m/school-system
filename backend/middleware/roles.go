package middleware

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func RequireRole(allowedRoles ...string) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			log.Printf("Проверка роли для запроса: %s %s", r.Method, r.URL.Path)
			log.Printf("Разрешенные роли: %v", allowedRoles)

			role, ok := r.Context().Value(ContextRole).(string)
			if !ok || role == "" {
				log.Printf("Роль не найдена в контексте")
				http.Error(w, "Нет роли в контексте", http.StatusUnauthorized)
				return
			}

			log.Printf("Текущая роль пользователя: %s", role)

			for _, allowed := range allowedRoles {
				if strings.EqualFold(role, allowed) {
					log.Printf("Доступ разрешен для роли: %s", role)
					next.ServeHTTP(w, r)
					return
				}
			}

			log.Printf("Доступ запрещен для роли: %s", role)
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(map[string]string{"error": "Недостаточно прав"})
		}
	}
}
