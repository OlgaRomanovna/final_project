package api

import (
	"net/http"
	"time"

	"todo-app/pkg/date"
)

func NextDateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSONError(w, http.StatusMethodNotAllowed, "метод не разрешен")
		return
	}

	dateStr := r.URL.Query().Get("date")
	repeat := r.URL.Query().Get("repeat")

	if dateStr == "" || repeat == "" {
		http.Error(w, "date and repeat are required", http.StatusBadRequest)
		return
	}

	// Проверяем, передан ли параметр now (для тестов)
	nowStr := r.URL.Query().Get("now")
	var now time.Time
	var err error

	if nowStr != "" {
		// Если передан, парсим его
		now, err = time.Parse("20060102", nowStr)
		if err != nil {
			http.Error(w, "invalid now format", http.StatusBadRequest)
			return
		}
	} else {
		// Иначе используем текущее время
		now = time.Now()
	}

	next, err := date.NextDate(now, dateStr, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	w.Write([]byte(next))
}
