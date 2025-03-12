package own_handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"
	"todo-app/internal/models"
	"todo-app/internal/repositories"
)

func GetTasksHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Устанавливаем заголовки CORS
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		userID := r.Context().Value("userID").(int)
		tasks, err := repositories.GetTasksByUserID(db, userID)
		if err != nil {
			log.Println("❌ Ошибка получения задач")
			http.Error(w, "Ошибка получения задач", http.StatusInternalServerError)
			return
		}

		jsonData, err := json.Marshal(tasks)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// Логируем успешное создание
		log.Printf("%s", "✅ Задачи получены")

		w.Write(jsonData)
	}
}

func CreateTaskHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Устанавливаем заголовки CORS
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		userID, ok := r.Context().Value("userID").(int)
		if !ok {
			log.Println("❌ Ошибка аутентификации: userID не найден")
			http.Error(w, "Ошибка аутентификации", http.StatusUnauthorized)
			return
		}

		// Временная структура для декодирования JSON
		var request struct {
			Title    string `json:"title"`
			Deadline string `json:"deadline"` // Получаем дату как строку
		}

		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			log.Printf("%s", "❌ Ошибка парсинга JSON: "+err.Error())
			http.Error(w, "Некорректный JSON", http.StatusBadRequest)
			return
		}

		// Проверяем, что поля не пустые
		if request.Title == "" || request.Deadline == "" {
			http.Error(w, "Название и дедлайн обязательны", http.StatusBadRequest)
			return
		}

		// Преобразуем строку в `time.Time`
		parsedDeadline, err := time.Parse("2006-01-02", request.Deadline)
		if err != nil {
			log.Println("❌ Ошибка парсинга даты:", err)
			http.Error(w, "Неверный формат даты. Используйте YYYY-MM-DD", http.StatusBadRequest)
			return
		}

		// Создаем задачу
		task := models.Task{
			UserID:   userID,
			Title:    request.Title,
			Deadline: parsedDeadline,
		}

		// Сохраняем в БД
		if err := repositories.CreateTask(db, &task); err != nil {
			log.Println("❌ Ошибка создания задачи в БД:", err)
			http.Error(w, "Ошибка создания задачи", http.StatusInternalServerError)
			return
		}
		// Логируем успешное создание
		log.Printf("✅ Задача создана: %s (до %s) пользователем %d", task.Title, task.Deadline.Format(time.RFC3339), userID)

		w.WriteHeader(http.StatusCreated)
	}
}

func DeleteTaskHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Устанавливаем заголовки CORS
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		userID, ok := r.Context().Value("userID").(int)
		if !ok {
			log.Println("❌ Ошибка аутентификации: userID не найден")
			http.Error(w, "Ошибка аутентификации", http.StatusUnauthorized)
			return
		}

		// Временная структура для декодирования JSON
		var request struct {
			Title    string `json:"title"`
			Deadline string `json:"deadline"` // Получаем дату как строку
		}

		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			log.Printf("%s", "❌ Ошибка парсинга JSON: "+err.Error())
			http.Error(w, "Некорректный JSON", http.StatusBadRequest)
			return
		}

		// Проверяем, что поля не пустые
		if request.Title == "" || request.Deadline == "" {
			log.Printf("%s", "❌ Название и дедлайн обязательны: title: "+request.Title+" deadline: "+request.Deadline)
			http.Error(w, "Название и дедлайн обязательны", http.StatusBadRequest)
			return
		}

		// Преобразуем строку в `time.Time`
		parsedDeadline, err := time.Parse("2006-01-02", request.Deadline)
		if err != nil {
			log.Println("❌ Ошибка парсинга даты:", err)
			http.Error(w, "Неверный формат даты. Используйте YYYY-MM-DD", http.StatusBadRequest)
			return
		}

		// Создаем задачу
		task := models.Task{
			UserID:   userID,
			Title:    request.Title,
			Deadline: parsedDeadline,
		}

		// Сохраняем в БД
		if err := repositories.DeleteTask(db, &task); err != nil {
			log.Println("❌ Ошибка удаления задачи в БД:", err)
			http.Error(w, "Ошибка удаления задачи", http.StatusInternalServerError)
			return
		}
		// Логируем успешное создание
		log.Printf("✅ Задача удалена: %s (до %s) пользователем %d", task.Title, task.Deadline.Format(time.RFC3339), userID)

		w.WriteHeader(http.StatusOK)
	}
}
