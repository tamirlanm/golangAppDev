package v1

import (
	"Practice7/internal/entity"
	"Practice7/internal/usecase"
	"Practice7/pkg/logger"
	"Practice7/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type userRoutes struct {
	t usecase.UserInterface
	l logger.Interface
}

func newUserRoutes(handler *gin.RouterGroup, t usecase.UserInterface, l logger.Interface) {
	r := &userRoutes{t, l}
	h := handler.Group("/users")
	{
		h.POST("/", r.RegisterUser)
		h.POST("/login", r.LoginUser)

		protected := h.Group("/")
		protected.Use(utils.JWTAuthMiddleware())
		{
			protected.GET("/protected/hello", r.ProtectedFunc)

			protected.GET("/me", r.GetMe)

			protected.PATCH("/promote/:id", utils.RoleMiddleware("admin"), r.PromoteUser)
		}
	}
}

func (r *userRoutes) RegisterUser(c *gin.Context) {
	var createUserDTO entity.CreateUserDTO
	if err := c.ShouldBindJSON(&createUserDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := utils.HashPassword(createUserDTO.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
		return
	}

	user := entity.User{
		Username: createUserDTO.Username,
		Email:    createUserDTO.Email,
		Password: hashedPassword,
		Role:     "user",
	}

	createdUser, sessionID, err := r.t.RegisterUser(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":    "User registered successfully. Please check your email for verification code.",
		"session_id": sessionID,
		"user":       createdUser,
	})
}

func (r *userRoutes) LoginUser(c *gin.Context) {
	var input entity.LoginUserDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := r.t.LoginUser(&input)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (r *userRoutes) ProtectedFunc(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "OK"})
}

func (r *userRoutes) GetMe(c *gin.Context) {
	userIDStr, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "userID not found in token"})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	user, err := r.t.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":       user.ID,
		"username": user.Username,
		"email":    user.Email,
		"role":     user.Role,
		"verified": user.Verified,
	})
}

func (r *userRoutes) PromoteUser(c *gin.Context) {
	idStr := c.Param("id")
	targetID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	if err := r.t.PromoteUser(targetID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User successfully promoted to admin",
		"user_id": targetID,
	})
}
