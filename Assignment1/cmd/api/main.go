package main

import (
	"Assignment1/internal/handlers"
	"Assignment1/internal/storage"
)

func main() {
	store := storage.NewTaskStorage()

	taskHandler := handlers.NewTaskHandler(store)

}
