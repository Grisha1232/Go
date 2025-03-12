package own_handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"todo-app/internal/services"
)

type Response struct {
	Message string `json:"message"`
}

// LoginHandler — обработчик входа пользователей
func LoginHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Устанавливаем заголовки CORS
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		var creds struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}

		if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
			log.Printf("❌ Некорректный запрос")
			http.Error(w, "Некорректный запрос", http.StatusBadRequest)
			return
		}

		token, err := services.Authenticate(db, creds.Username, creds.Password)
		if err != nil {
			log.Printf("%s", "❌ Ошибка входа: "+err.Error())
			http.Error(w, "Ошибка входа: "+err.Error(), http.StatusUnauthorized)
			return
		}

		// Отправляем JSON-ответ с токеном
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		jsonData, err := json.Marshal(map[string]string{"token": token})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Printf("%s", "✅ Вход успешен "+creds.Username+" "+token)
		w.Write(jsonData)
	}
}

// RegisterHandler — обработчик регистрации пользователей
func RegisterHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Устанавливаем заголовки CORS
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		var creds struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}

		if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
			log.Printf("❌ Некорректный запрос ")
			http.Error(w, "Некорректный запрос", http.StatusBadRequest)
			return
		}

		err := services.RegisterUser(db, creds.Username, creds.Password)
		if err != nil {
			if err.Error() == "пользователь уже существует" {
				log.Printf("❌ Пользователь уже существует ")
				http.Error(w, err.Error(), http.StatusConflict) // 409 Conflict
			} else {
				log.Printf("%s", "❌ Ошибка регистрации: "+err.Error())
				http.Error(w, "Ошибка регистрации: "+err.Error(), http.StatusInternalServerError)
			}
			return
		}

		// Отправляем JSON-ответ с кодом 201 Created
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		jsonData, err := json.Marshal(map[string]string{"message": "Регистрация успешна"})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Printf("%s", "✅ Регистрация успешена "+creds.Username)
		w.Write(jsonData)
	}
}
