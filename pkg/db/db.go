package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

const schema = `
CREATE TABLE IF NOT EXISTS scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT '',
    title VARCHAR(255) NOT NULL DEFAULT '',
    comment TEXT,
    repeat VARCHAR(128) NOT NULL DEFAULT ''
);
CREATE INDEX IF NOT EXISTS idx_scheduler_date ON scheduler(date);
`

func Init(dbFile string) error {
	_, err := os.Stat(dbFile)
	install := err != nil

	DB, err = sql.Open("sqlite", dbFile)
	if err != nil {
		return fmt.Errorf("ошибка открытия базы данных: %w", err)
	}

	if err := DB.Ping(); err != nil {
		DB.Close() // <-- Закрываем при ошибке
		return fmt.Errorf("ошибка подключения к базе данных: %w", err)
	}

	if install {
		if _, err := DB.Exec(schema); err != nil {
			DB.Close() // <-- Закрываем при ошибке
			return fmt.Errorf("ошибка выполнения SQL-скрипта инициализации: %w", err)
		}
	}
	return nil
}
