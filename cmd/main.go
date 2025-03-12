package main

import (
	"fmt"
	"log"
	"net/http"

	"todo-app/config"
	"todo-app/internal/database"
	"todo-app/internal/handlers"
	"todo-app/internal/middleware"

	"github.com/gorilla/mux"
)

func main() {
	// Загружаем конфигурацию
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("❌ Ошибка загрузки конфигурации: %v", err)
	}

	// Подключаемся к базе данных
	db, err := database.ConnectDB(cfg)
	if err != nil {
		log.Fatalf("❌ Ошибка подключения к БД: %v", err)
	}
	defer db.Close()

	// Запускаем миграции
	if err := database.RunMigrations(db); err != nil {
		log.Fatalf("❌ Ошибка выполнения миграций: %v", err)
	}

	// Создаем роутер
	r := mux.NewRouter()

	// ВАЖНО: сначала подключаем CORS Middleware
	r.Use(middleware.CORSMiddleware)

	// Логирование запросов
	r.Use(middleware.LoggingMiddleware)

	// Эндпоинты аутентификации
	r.HandleFunc("/login", handlers.LoginHandler(db)).Methods("POST", "OPTIONS")

	// Эндпоинты задач (доступны только аутентифицированным пользователям)
	r.HandleFunc("/tasks", middleware.AuthMiddleware(handlers.GetTasksHandler(db))).Methods("GET", "OPTIONS")
	r.HandleFunc("/tasks", middleware.AuthMiddleware(handlers.CreateTaskHandler(db))).Methods("POST", "OPTIONS")

	// Запуск сервера
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("🚀 Сервер запущен на http://localhost%s", addr)
	log.Fatal(http.ListenAndServe(addr, r))
}
