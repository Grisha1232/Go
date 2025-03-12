package middleware

import (
	"log"
	"net/http"
)

// ErrorLoggingMiddleware логирует ошибки при обработке запросов
func ErrorLoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Создаем ResponseWriter, который будет отслеживать статус-код ответа
		lrw := &loggingResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(lrw, r)

		// Если статус-код 4xx или 5xx, логируем ошибку
		if lrw.statusCode >= 400 {
			log.Printf("❌ Ошибка HTTP %d при запросе %s %s", lrw.statusCode, r.Method, r.URL.Path)
		}
	})
}

// loggingResponseWriter — обертка для ResponseWriter, чтобы отслеживать статус ответа
type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}
