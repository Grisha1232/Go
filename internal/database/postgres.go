package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"todo-app/config"

	_ "github.com/lib/pq"
)

// ConnectDB подключается к PostgreSQL и возвращает соединение
func ConnectDB(cfg *config.Config) (*sql.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Database.Host, cfg.Database.Port, cfg.Database.User, cfg.Database.Password, cfg.Database.DBName)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	log.Println("✅ Успешное подключение к базе данных")
	return db, nil
}

// RunMigrations автоматически применяет миграции перед запуском приложения
func RunMigrations(db *sql.DB) error {
	log.Println("🚀 Проверка и запуск миграций...")

	// Проверяем, существует ли таблица "users"
	var exists bool
	err := db.QueryRow("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'users')").Scan(&exists)
	if err != nil {
		return fmt.Errorf("ошибка проверки существования таблицы users: %v", err)
	}

	if exists {
		log.Println("✅ Таблица users уже существует, пропускаем миграцию.")
		return nil // Если таблица есть, миграция не требуется
	}

	// Читаем SQL-миграцию
	migrationFile := "internal/database/migrations/001_init.up.sql"
	migration, err := os.ReadFile(migrationFile)
	if err != nil {
		return fmt.Errorf("ошибка чтения файла миграции %s: %v", migrationFile, err)
	}

	// Выполняем миграцию
	_, err = db.Exec(string(migration))
	if err != nil {
		return fmt.Errorf("ошибка выполнения миграции %s: %v", migrationFile, err)
	}

	log.Printf("✅ Миграция %s успешно применена", migrationFile)
	return nil
}
