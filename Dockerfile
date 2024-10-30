# Используем образ Go с нужной версией
FROM golang:1.21-alpine AS builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем файлы go.mod и go.sum перед тем, как копировать весь исходный код.
# Это позволит Docker использовать кеш, если зависимости не изменились.
COPY go.mod go.sum ./
RUN go mod tidy
RUN go mod download

# Копируем весь исходный код после установки зависимостей
COPY . .

# Собираем приложение
RUN go build -o app main.go

# Используем более легкий базовый образ для финального контейнера
FROM alpine:latest

# Копируем только скомпилированный бинарный файл и .env
COPY --from=builder /app/app /app/app
COPY --from=builder /app/.env /app/.env

# Устанавливаем рабочую директорию
WORKDIR /app

# Указываем точку входа с аргументом, который будет передан при запуске контейнера
ENTRYPOINT ["/app/app"]
