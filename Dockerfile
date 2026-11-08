# Этап 1: Сборка приложения
FROM golang:1.26-alpine AS builder

WORKDIR /app

# Копируем файлы модулей и скачиваем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем статический бинарный файл для Linux
RUN CGO_ENABLED=0 GOOS=linux go build -o scheduler .

# Этап 2: Финальный минимальный образ
FROM alpine:latest

WORKDIR /app

# Копируем только скомпилированный бинарник и папку с фронтендом
COPY --from=builder /app/scheduler .
COPY --from=builder /app/web ./web

# Указываем порт, который будет слушать приложение
EXPOSE 7540

# Переменные окружения по умолчанию (можно переопределить при запуске)
ENV TODO_PORT=7540
ENV TODO_DBFILE=/app/scheduler.db

# Команда запуска
CMD ["./scheduler"]