package storage

import (
	"Assignment1/internal/models"
	"errors"
	"sync"
)

type TaskStorage struct {
	tasks  map[int]models.Task
	nextID int
	mu     sync.RWMutex
}

func NewTaskStorage() *TaskStorage {
	return &TaskStorage{
		tasks:  make(map[int]models.Task),
		nextID: 1,
	}
}

func (s *TaskStorage) Create(title string) models.Task {
	s.mu.Lock()
	defer s.mu.Unlock()

	task := models.Task{
		ID:    s.nextID,
		Title: title,
		Done:  false,
	}
	s.tasks[task.ID] = task
	s.nextID++

	return task
}

func (s *TaskStorage) GetAll() []models.Task {
	s.mu.RLock()
	defer s.mu.RUnlock()

	res := make([]models.Task, 0, len(s.tasks))

	for _, task := range s.tasks {
		res = append(res, task)
	}
	return res
}

func (s *TaskStorage) GetByID(id int) (models.Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	task, exists := s.tasks[id]
	if !exists {
		return models.Task{}, errors.New("task not found")
	}
	return task, nil
}

func (s *TaskStorage) Update(id int, done bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	task, exists := s.tasks[id]
	if !exists {
		return errors.New("task not found")
	}

	task.Done = done
	s.tasks[id] = task
	return nil
}
