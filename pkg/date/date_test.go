package date

import (
	"testing"
	"time"
)

func TestNextDate(t *testing.T) {
	// Фиксируем "сегодня" как 26.01.2024
	now := time.Date(2024, 1, 26, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name        string
		now         time.Time
		dstart      string
		repeat      string
		expected    string
		expectError bool
	}{
		// Базовые правила
		{"d 7 - пример из задания", now, "20240113", "d 7", "20240127", false},
		{"d 1 - каждый день", now, "20240125", "d 1", "20240127", false},
		{"y - ежегодно (високосный год)", now, "20240229", "y", "20250301", false},
		{"y - ежегодно (обычный год)", now, "20230315", "y", "20240315", false},

		// Правила со звёздочкой
		{"w 7 - ближайшее воскресенье", now, "20240120", "w 7", "20240128", false},
		{"w 1,4,5 - пн, чт, пт", now, "20240120", "w 1,4,5", "20240129", false},
		{"m 16,5 - пример из задания", now, "20240116", "m 16,5", "20240205", false},
		{"m -1,18 - пример из задания", now, "20240201", "m -1,18", "20240218", false},
		{"m 4 - 4-е число каждого месяца", now, "20240101", "m 4", "20240204", false},
		{"m 3 1,3,6 - 3 января, марта, июня", now, "20240101", "m 3 1,3,6", "20240303", false},

		// Ошибки
		{"Пустое правило", now, "20240101", "", "", true},
		{"Неверная дата dstart", now, "invalid", "d 1", "", true},
		{"Недопустимый символ", now, "20240101", "l 34", "", true},
		{"d без интервала", now, "20240101", "d", "", true},
		{"d 405 - превышен интервал", now, "20240101", "d 405", "", true},
		{"w 8 - недопустимый день", now, "20240101", "w 8", "", true},
		{"m 4,40 - недопустимый день месяца", now, "20240101", "m 4,40", "", true},
		{"m 5 1,13 - недопустимый месяц", now, "20240101", "m 5 1,13", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := NextDate(tt.now, tt.dstart, tt.repeat)
			if tt.expectError {
				if err == nil {
					t.Errorf("ожидалась ошибка, но её не было. Результат: %s", result)
				}
			} else {
				if err != nil {
					t.Errorf("неожиданная ошибка: %v", err)
				}
				if result != tt.expected {
					t.Errorf("ожидалось %s, получено %s", tt.expected, result)
				}
			}
		})
	}
}
