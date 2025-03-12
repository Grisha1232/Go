# Используем официальный образ Go для сборки
FROM golang:1.24 AS builder

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем Go-модули и загружаем зависимости
COPY go.mod go.sum ./

# Копируем весь исходный код
COPY . .

# Компилируем бинарник для Linux
RUN GOOS=linux go build -o todo-app cmd/app/main.go

# Минимальный продакшн-образ
FROM alpine:latest as production

RUN addgroup -S appgroup && adduser -S appuser -G appgroup

USER appuser

# Устанавливаем зависимости (например, ca-certificates)
RUN apk --no-cache add ca-certificates

# Копируем скомпилированное приложение
COPY --from=build /app/todo-app ./

# Делаем бинарник исполняемым
RUN chmod +x ./todo-app

# Открываем порт 8080
EXPOSE 8080

# Запускаем приложение
CMD ["./todo-app"]
