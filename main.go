package main

import (
	"log"
	"os"
	"todo-app/pkg/db"
)

func main() {
	cfg := LoadConfig()
	log.Printf("Инициализация приложения. WebDir: %s, Port: %d", cfg.WebDir, cfg.Port)

	// Читаем путь к БД из переменной окружения, если она задана тестами
	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		dbFile = "scheduler.db"
	}

	log.Printf("Проверка и инициализация базы данных: %s", dbFile)
	if err := db.Init(dbFile); err != nil {
		log.Fatalf("Критическая ошибка инициализации БД: %v", err)
	}
	log.Println("База данных успешно инициализирована")

	if err := StartServer(cfg); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
