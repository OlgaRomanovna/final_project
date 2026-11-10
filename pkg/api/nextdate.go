package api

import (
	"log"
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

	nowStr := r.URL.Query().Get("now")
	var now time.Time
	var err error

	if nowStr != "" {
		now, err = time.Parse("20060102", nowStr)
		if err != nil {
			http.Error(w, "invalid now format", http.StatusBadRequest)
			return
		}
	} else {
		now = time.Now()
	}

	next, err := date.NextDate(now, dateStr, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	if _, err := w.Write([]byte(next)); err != nil {
		log.Printf("ошибка записи ответа: %v", err)
	}
}
