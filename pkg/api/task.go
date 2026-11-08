package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"todo-app/pkg/date"
	"todo-app/pkg/db"
)

// TaskHandler маршрутизирует запросы к /api/task в зависимости от HTTP-метода
func TaskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		addTaskHandler(w, r)
	case http.MethodGet:
		getTaskHandler(w, r)
	case http.MethodPut:
		updateTaskHandler(w, r)
	case http.MethodDelete:
		deleteTaskHandler(w, r)
	default:
		writeJSONError(w, http.StatusMethodNotAllowed, "метод не разрешен")
	}
}

// TaskDoneHandler обрабатывает отметку задачи как выполненной (POST /api/task/done)
func TaskDoneHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSONError(w, http.StatusMethodNotAllowed, "метод не разрешен")
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSONError(w, http.StatusBadRequest, "не указан идентификатор")
		return
	}

	// 1. Получаем текущие параметры задачи
	task, err := db.GetTask(id)
	if err != nil {
		writeJSONError(w, http.StatusNotFound, err.Error())
		return
	}

	// 2. Логика в зависимости от наличия правила повторения
	if task.Repeat == "" {
		// Одноразовая задача: удаляем
		err = db.DeleteTask(id)
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, err.Error())
			return
		}
	} else {
		// Периодическая задача: вычисляем новую дату и обновляем
		now := time.Now()
		now = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

		nextDate, err := date.NextDate(now, task.Date, task.Repeat)
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, "ошибка вычисления следующей даты: "+err.Error())
			return
		}

		err = db.UpdateDate(nextDate, id)
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	// 3. Возвращаем пустой JSON объект {}
	writeJSON(w, struct{}{})
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeJSONError(w, http.StatusBadRequest, "ошибка десериализации JSON")
		return
	}

	if task.Title == "" {
		writeJSONError(w, http.StatusBadRequest, "не указан заголовок задачи")
		return
	}

	if err := checkAndFixDate(&task); err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	id, err := db.AddTask(&task)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "ошибка сохранения в базу данных")
		return
	}

	writeJSON(w, map[string]string{"id": fmt.Sprintf("%d", id)})
}

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSONError(w, http.StatusBadRequest, "не указан идентификатор")
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeJSONError(w, http.StatusNotFound, err.Error())
		return
	}

	writeJSON(w, task)
}

func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeJSONError(w, http.StatusBadRequest, "ошибка десериализации JSON")
		return
	}

	if task.ID == "" {
		writeJSONError(w, http.StatusBadRequest, "не указан идентификатор задачи")
		return
	}

	if task.Title == "" {
		writeJSONError(w, http.StatusBadRequest, "не указан заголовок задачи")
		return
	}

	// Проверяем и корректируем дату по тем же правилам, что и при создании
	if err := checkAndFixDate(&task); err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	err := db.UpdateTask(&task)
	if err != nil {
		writeJSONError(w, http.StatusNotFound, err.Error())
		return
	}

	// Возвращаем пустой JSON объект {}
	writeJSON(w, struct{}{})
}

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSONError(w, http.StatusBadRequest, "не указан идентификатор")
		return
	}

	err := db.DeleteTask(id)
	if err != nil {
		writeJSONError(w, http.StatusNotFound, err.Error())
		return
	}

	writeJSON(w, struct{}{})
}

func checkAndFixDate(task *db.Task) error {
	now := time.Now()
	now = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	if task.Date == "" {
		task.Date = now.Format("20060102")
	}

	t, err := time.Parse("20060102", task.Date)
	if err != nil {
		return fmt.Errorf("неверный формат даты")
	}

	if task.Repeat != "" {
		next, err := date.NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return fmt.Errorf("неверный формат правила повторения")
		}

		// Двигаем только если дата строго меньше сегодняшней
		if now.After(t) {
			task.Date = next
		}
	} else {
		// Ставим сегодня только если дата строго меньше сегодняшней
		if now.After(t) {
			task.Date = now.Format("20060102")
		}
	}

	return nil
}

func writeJSON(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(data)
}

func writeJSONError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}
