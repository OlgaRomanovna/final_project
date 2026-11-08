package main

import (
	"os"
	"strconv"
)

const (
	DefaultPort   = 7540
	DefaultWebDir = "./web"
)

// Config хранит настройки приложения
type Config struct {
	Port   int
	WebDir string
}

// LoadConfig загружает конфигурацию из переменных окружения или использует значения по умолчанию
func LoadConfig() Config {
	port := DefaultPort

	// Проверяем, задана ли переменная окружения TODO_PORT
	if portEnv := os.Getenv("TODO_PORT"); portEnv != "" {
		if p, err := strconv.Atoi(portEnv); err == nil {
			port = p
		}
	}

	return Config{
		Port:   port,
		WebDir: DefaultWebDir,
	}
}
