package services

import (
	"database/sql"
	"errors"
	"todo-app/internal/models"
	"todo-app/internal/repositories"
)

// GetUserTasks возвращает список задач пользователя
func GetUserTasks(db *sql.DB, userID int) ([]models.Task, error) {
	tasks, err := repositories.GetTasksByUserID(db, userID)
	if err != nil {
		return nil, errors.New("ошибка получения задач")
	}
	return tasks, nil
}

// AddUserTask добавляет новую задачу в БД
func AddUserTask(db *sql.DB, task *models.Task) error {
	if task.Title == "" || task.Deadline.IsZero() {
		return errors.New("название и дедлайн обязательны")
	}

	err := repositories.CreateTask(db, task)
	if err != nil {
		return errors.New("ошибка создания задачи")
	}

	return nil
}

func DeleteUserTask(db *sql.DB, task *models.Task) error {
	if task.Title == "" || task.Deadline.IsZero() {
		return errors.New("название и дедлайн обязательны")
	}

	err := repositories.DeleteTask(db, task)
	if err != nil {
		return errors.New("ошибка удаления задачи")
	}

	return nil
}
