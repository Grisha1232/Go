package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"todo-app/internal/models"
	"todo-app/internal/repositories"
)

func GetTasksHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value("userID").(int)
		tasks, err := repositories.GetTasksByUserID(db, userID)
		if err != nil {
			http.Error(w, "Ошибка получения задач", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(tasks)
	}
}

func CreateTaskHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var task models.Task
		if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
			http.Error(w, "Некорректный запрос", http.StatusBadRequest)
			return
		}

		userID := r.Context().Value("userID").(int)
		task.UserID = userID

		if err := repositories.CreateTask(db, &task); err != nil {
			http.Error(w, "Ошибка создания задачи", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
	}
}
