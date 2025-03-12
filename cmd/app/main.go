package main

import (
	"fmt"
	"log"
	"net/http"

	"todo-app/config"
	"todo-app/internal/database"
	"todo-app/internal/middleware"
	"todo-app/internal/own_handlers"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Метрики Prometheus
var (
	requestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Общее количество HTTP-запросов",
		},
		[]string{"method", "path"},
	)

	requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Длительность HTTP-запросов",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)
)

// Middleware для сбора метрик
func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		timer := prometheus.NewTimer(requestDuration.WithLabelValues(r.Method, r.URL.Path))
		defer timer.ObserveDuration()
		requestsTotal.WithLabelValues(r.Method, r.URL.Path).Inc()
		next.ServeHTTP(w, r)
	})
}

func main() {
	// Регистрация метрик
	prometheus.MustRegister(requestsTotal, requestDuration)

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

	// ✅ Middleware для CORS (разрешает frontend доступ к backend)
	corsMiddleware := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),                                // Разрешить все домены (изменить в продакшене)
		handlers.AllowedMethods([]string{"GET", "DELETE", "POST", "OPTIONS"}), // Разрешенные методы
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),    // Разрешенные заголовки
	)

	// Логирование запросов
	r.Use(middleware.LoggingMiddleware)
	r.Use(middleware.ErrorLoggingMiddleware)
	r.Use(MetricsMiddleware) // Подключаем метрики

	// ✅ Добавляем обработчик OPTIONS-запросов (preflight CORS)
	r.Methods("OPTIONS").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Роут для регистрации пользователя
	r.HandleFunc("/register", own_handlers.RegisterHandler(db)).Methods("POST", "OPTIONS")

	// Эндпоинты аутентификации
	r.HandleFunc("/login", own_handlers.LoginHandler(db)).Methods("POST", "OPTIONS")

	// Эндпоинты задач (доступны только аутентифицированным пользователям)
	r.HandleFunc("/tasks", middleware.AuthMiddleware(own_handlers.GetTasksHandler(db))).Methods("GET", "OPTIONS")
	r.HandleFunc("/tasks", middleware.AuthMiddleware(own_handlers.CreateTaskHandler(db))).Methods("POST", "OPTIONS")
	r.HandleFunc("/tasks", middleware.AuthMiddleware(own_handlers.DeleteTaskHandler(db))).Methods("DELETE", "OPTIONS")

	// Эндпоинт для Prometheus-метрик
	r.Handle("/metrics", promhttp.Handler())

	// Запуск сервера
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("🚀 Сервер запущен на http://localhost%s", addr)
	log.Fatal(http.ListenAndServe(addr, corsMiddleware(r)))
	log.Fatal(http.ListenAndServe(addr, r))
}
