package main

import (
	"Assignment1/internal/handlers"
	"Assignment1/internal/middleware"
	"Assignment1/internal/storage"
	"log"
	"net/http"
)

func main() {
	store := storage.NewTaskStorage()

	taskHandler := handlers.NewTaskHandler(store)
	mux := http.NewServeMux()
	mux.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			taskHandler.GetTasks(w, r)
		case http.MethodPost:
			taskHandler.CreateTask(w, r)
		case http.MethodPatch:
			taskHandler.UpdateTask(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	handler := middleware.Logging(middleware.APIKeyAuth(mux))
	log.Println("Server is running on :8080")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatal(err)
	}
}
