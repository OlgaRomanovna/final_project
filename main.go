package main

import (
	"log"

	"todo-app/pkg/api"
	"todo-app/pkg/db"
)

func main() {
	cfg := LoadConfig()
	log.Printf("Инициализация приложения. WebDir: %s, Port: %d, DBFile: %s", cfg.WebDir, cfg.Port, cfg.DBFile)

	api.Init(cfg.Password)

	if err := db.Init(cfg.DBFile); err != nil {
		log.Fatalf("Критическая ошибка инициализации БД: %v", err)
	}
	defer db.DB.Close()
	log.Println("База данных успешно инициализирована")

	if err := StartServer(cfg); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
