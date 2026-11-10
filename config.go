package main

import (
	"os"
	"strconv"
)

const (
	DefaultPort   = 7540
	DefaultWebDir = "./web"
)

type Config struct {
	Port     int
	WebDir   string
	DBFile   string
	Password string
}

func LoadConfig() Config {
	port := DefaultPort
	if p, err := strconv.Atoi(os.Getenv("TODO_PORT")); err == nil {
		port = p
	}

	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		dbFile = "scheduler.db"
	}

	return Config{
		Port:     port,
		WebDir:   DefaultWebDir,
		DBFile:   dbFile,
		Password: os.Getenv("TODO_PASSWORD"),
	}
}
