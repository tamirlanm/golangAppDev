package app

import (
	"context"
	"log"
	"net/http"
	"time"

	middleware "Practice3/internal/handler"
	userHandler "Practice3/internal/handler/users"
	"Practice3/internal/repository"
	_postgres "Practice3/internal/repository/_postgres"
	usecase "Practice3/internal/usecase/users"
	"Practice3/pkg/modules"
)

func Run() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dbConfig := initPostgresConfig()
	db := _postgres.NewPGXDialect(ctx, dbConfig)

	repos := repository.NewRepositories(db)
	uc := usecase.NewUserUsecase(repos.UserRepository)
	uh := userHandler.NewUserHandler(uc)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /healthcheck", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"OK"}`))
	})

	mux.HandleFunc("GET /users", uh.GetUsers)
	mux.HandleFunc("GET /users/{id}", uh.GetUserByID)
	mux.HandleFunc("POST /users", uh.CreateUser)
	mux.HandleFunc("PUT /users/{id}", uh.UpdateUser)
	mux.HandleFunc("DELETE /users/{id}", uh.DeleteUser)

	chain := middleware.LoggingMiddleware(middleware.AuthMiddleware(mux))
	log.Println("Server is running on port 8080")
	if err := http.ListenAndServe(":8080", chain); err != nil {
		log.Fatal(err)
	}

}

func initPostgresConfig() *modules.PostgreConfig {
	return &modules.PostgreConfig{
		Host:        "localhost",
		Port:        "5432",
		Username:    "postgres",
		Password:    "12345",
		DBName:      "mygodb",
		SSLMode:     "disable",
		ExecTimeout: 5 * time.Second,
	}
}
