package handlers

import (
	"Assignment1/internal/storage"
)

type TaskHandler struct {
	storage *storage.TaskStorage
}

func NewTaskHandler(storage *storage.TaskStorage) *TaskHandler {
	return &TaskHandler{storage: storage}
}
