package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "modernc.org/sqlite"
)

// DB - глобальная переменная для хранения подключения к базе данных.
var DB *sql.DB

// schema содержит SQL-команды для создания таблицы и индекса
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

// Init открывает базу данных и создает таблицу с индексом, если файла не существует
func Init(dbFile string) error {
	// Проверяем наличие файла базы данных
	_, err := os.Stat(dbFile)

	var install bool
	if err != nil {
		// Если файл не найден (или другая ошибка доступа), помечаем для установки
		install = true
	}

	// Открываем соединение с базой данных
	DB, err = sql.Open("sqlite", dbFile)
	if err != nil {
		return fmt.Errorf("ошибка открытия базы данных: %w", err)
	}

	// Проверяем, что соединение действительно работает
	if err := DB.Ping(); err != nil {
		return fmt.Errorf("ошибка подключения к базе данных: %w", err)
	}

	// Если файла не было, выполняем SQL-скрипт создания таблицы и индекса
	if install {
		if _, err := DB.Exec(schema); err != nil {
			return fmt.Errorf("ошибка выполнения SQL-скрипта инициализации: %w", err)
		}
	}

	return nil
}
