package v1

import (
	"Practice7/internal/usecase"
	"Practice7/pkg/logger"
	"Practice7/utils"

	"github.com/gin-gonic/gin"
)

func NewRouter(handler *gin.Engine, t usecase.UserInterface, l logger.Interface) {
	handler.Use(utils.RateLimiterMiddleware())

	v1 := handler.Group("/v1")
	{
		newUserRoutes(v1, t, l)
	}
}
