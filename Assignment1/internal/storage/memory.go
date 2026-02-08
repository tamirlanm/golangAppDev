package storage

import (
	"Assignment1/internal/handlers"
	"errors"
	"sync"
)

type TaskStorage struct {
	tasks  map[int]handlers.Task
	nextID int
	mu     sync.RWMutex
}

func NewTaskStorage() *TaskStorage {
	return &TaskStorage{
		tasks:  make(map[int]handlers.Task),
		nextID: 1,
	}
}

func (s *TaskStorage) Create(title string) handlers.Task {
	s.mu.Lock()
	defer s.mu.Unlock()

	task := handlers.Task{
		ID:    s.nextID,
		Title: title,
		Done:  false,
	}
	s.tasks[task.ID] = task
	s.nextID++

	return task
}

func (s *TaskStorage) GetAll() []handlers.Task {
	s.mu.RLock()
	defer s.mu.RUnlock()

	res := make([]handlers.Task, 0, len(s.tasks))

	for _, task := range s.tasks {
		res = append(res, task)
	}
	return res
}

func (s *TaskStorage) GetByID(id int) (handlers.Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	task, exists := s.tasks[id]
	if !exists {
		return handlers.Task{}, errors.New("task not found")
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
