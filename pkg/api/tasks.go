package api

import (
	"net/http"

	"todo-app/pkg/db"
)

const tasksLimit = 50

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func TasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSONError(w, http.StatusMethodNotAllowed, "метод не разрешен")
		return
	}

	tasks, err := db.Tasks(tasksLimit)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, TasksResp{
		Tasks: tasks,
	})
}
