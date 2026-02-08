package handlers

import (
	"Assignment1/internal/models"
	"Assignment1/internal/storage"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

type TaskHandler struct {
	storage *storage.TaskStorage
}

func NewTaskHandler(storage *storage.TaskStorage) *TaskHandler {
	return &TaskHandler{storage: storage}
}

func (h *TaskHandler) GetTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	idStr := r.URL.Query().Get("id")

	if idStr != "" {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "invalid id",
			})
			return
		}
		task, err := h.storage.GetByID(id)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "taks not found",
			})
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(task)
		return
	}

	tasks := h.storage.GetAll()
	doneStr := r.URL.Query().Get("done")
	if doneStr != "" {
		doneValue := doneStr == "true"

		filtered := make([]models.Task, 0)
		for _, task := range tasks {
			if task.Done == doneValue {
				filtered = append(filtered, task)
			}
		}
		tasks = filtered
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tasks)
}

func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var req struct {
		Title string `json:"title"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "invalid request body",
		})
		return
	}

	if strings.TrimSpace(req.Title) == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "invalid title",
		})
		return
	}
	if len(req.Title) > 150 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "title too long (150 characters max)",
		})
		return
	}

	task := h.storage.Create(req.Title)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)

}

func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "id is required",
		})
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "invalid id",
		})
		return
	}

	var req struct {
		Done bool `json:"done"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "invalid request body",
		})
		return
	}

	if err := h.storage.Update(id, req.Done); err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "task not found",
		})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]bool{
		"updated": true,
	})
}
