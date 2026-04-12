package app

import (
	"Practice7/config"
	v1 "Practice7/internal/controller/http/v1"
	"Practice7/internal/entity"
	"Practice7/internal/usecase"
	"Practice7/internal/usecase/repo"
	"Practice7/pkg/logger"
	"Practice7/pkg/postgres"
	"log"

	"github.com/gin-gonic/gin"
)

func Run(cfg *config.Config) {
	pg, err := postgres.New(cfg.DBURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	pg.Conn.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"")
	if err := pg.Conn.AutoMigrate(&entity.User{}); err != nil {
		log.Fatal("Failed to migrate:", err)
	}

	l := logger.New()
	userRepo := repo.NewUserRepo(pg)
	userUC := usecase.NewUserUseCase(userRepo)
	router := gin.Default()
	v1.NewRouter(router, userUC, l)

	port := cfg.Port
	if port == "" {
		port = "8090"
	}
	log.Println("Server started on port", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to run server:", err)
	}
}
