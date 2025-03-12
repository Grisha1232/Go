# Используем официальный образ Go для сборки
FROM golang:1.24 AS build

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем весь исходный код
COPY . .

ENV CGO_ENABLED=0
ENV GOOS=linux

# Компилируем бинарник для Linux
RUN go build -o todo-app cmd/app/main.go

# Минимальный продакшн-образ
FROM alpine:latest as production

RUN addgroup -S appgroup && adduser -S appuser -G appgroup

USER appuser

# Копируем скомпилированное приложение
COPY --from=build /app/todo-app ./
COPY --from=build /app/config ./config
COPY --from=build /app/internal/ ./internal

# Открываем порт 8080
EXPOSE 8080

# Запускаем приложение
CMD ["./todo-app"]
