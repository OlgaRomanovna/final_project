package main

import (
	"fmt"
	"log"
	"net/http"

	"todo-app/pkg/api"
)

func StartServer(cfg Config) error {
	fileServer := http.FileServer(http.Dir(cfg.WebDir))
	http.Handle("/", fileServer)

	// Публичный маршрут для входа
	http.HandleFunc("/api/signin", api.SignInHandler)
	http.HandleFunc("/api/nextdate", api.NextDateHandler)

	// Защищенные маршруты (оборачиваем в AuthMiddleware)
	http.HandleFunc("/api/task", api.AuthMiddleware(api.TaskHandler))
	http.HandleFunc("/api/tasks", api.AuthMiddleware(api.TasksHandler))
	http.HandleFunc("/api/task/done", api.AuthMiddleware(api.TaskDoneHandler))

	addr := fmt.Sprintf(":%d", cfg.Port)
	log.Printf("Сервер запущен и слушает порт %d", cfg.Port)
	log.Printf("Откройте в браузере: http://localhost%s", addr)

	return http.ListenAndServe(addr, nil)
}
