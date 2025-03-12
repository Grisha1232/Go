# Используем официальный образ Go для сборки
FROM golang:1.18 AS builder

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем Go-модули и загружаем зависимости
COPY go.mod go.sum ./
RUN go mod tidy

# Копируем весь исходный код
COPY . .

# Собираем бинарник
RUN go build -o todo-app cmd/app/main.go

# Создаем минимальный образ для продакшена
FROM alpine:latest

# Устанавливаем зависимости (например, ca-certificates)
RUN apk --no-cache add ca-certificates

# Устанавливаем рабочую директорию
WORKDIR /root/

# Копируем бинарник из builder-слоя
COPY --from=builder /app/todo-app .

# Открываем порт 8080
EXPOSE 8080

# Запускаем приложение
CMD ["./todo-app"]
