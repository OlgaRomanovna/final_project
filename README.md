# Планировщик задач (TODO-лист)

Веб-сервер на Go для управления задачами с поддержкой правил повторения (ежедневно, еженедельно, ежемесячно, ежегодно). Данные хранятся в SQLite. Реализует REST API, с которым работает готовый фронтенд.

## Запуск проекта

Сборка образа:
```bash
docker build -t todo-scheduler .
```

```bash
docker run -d \
--name scheduler-app \
-p 7540:7540 \
-v $(pwd)/scheduler.db:/app/scheduler.db \
-e TODO_PASSWORD=12345 \
todo-scheduler
```

Откройте в браузере: http://localhost:7540/login.html и войдите с паролем 12345

## Запуск тестов
Перед запуском тестов получите JWT-токен через POST-запрос к /api/signin с телом {"password": "12345"} и вставьте его в переменную Token в файле tests/settings.go

```bash
go test -v ./tests
```