package repositories

import (
	"database/sql"
	"todo-app/internal/models"
)

func GetTasksByUserID(db *sql.DB, userID int) ([]models.Task, error) {
	rows, err := db.Query("SELECT id, title, description, deadline, completed FROM tasks WHERE user_id=$1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var task models.Task
		if err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.Deadline, &task.Completed); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func CreateTask(db *sql.DB, task *models.Task) error {
	_, err := db.Exec("INSERT INTO tasks (user_id, title, description, deadline, completed) VALUES ($1, $2, $3, $4, $5)",
		task.UserID, task.Title, task.Description, task.Deadline, task.Completed)
	return err
}
