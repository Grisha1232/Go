package repositories

import (
	"database/sql"
	"todo-app/internal/models"
)

func GetTasksByUserID(db *sql.DB, userID int) ([]models.Task, error) {
	rows, err := db.Query("SELECT id, user_id, title, deadline FROM tasks WHERE user_id=$1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var task models.Task
		if err := rows.Scan(&task.ID, &task.UserID, &task.Title, &task.Deadline); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

// CreateTask создает новую задачу в БД
func CreateTask(db *sql.DB, task *models.Task) error {
	query := `INSERT INTO tasks (user_id, title, deadline) VALUES ($1, $2, $3)`
	_, err := db.Exec(query, task.UserID, task.Title, task.Deadline)
	return err
}

func DeleteTask(db *sql.DB, task *models.Task) error {
	query := `DELETE FROM tasks WHERE (userd_id = $1 AND title = $2 AND deadline = $3)`
	_, err := db.Exec(query, task.UserID, task.Title, task.Deadline)
	return err
}
