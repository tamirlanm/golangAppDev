package utils

import (
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type clientData struct {
	count       int
	windowStart time.Time
}

var (
	rateMu  sync.Mutex
	clients = make(map[string]*clientData)
)

const (
	rateLimit  = 5
	rateWindow = time.Minute
)

func RateLimiterMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		key := resolveClientKey(c)

		rateMu.Lock()
		defer rateMu.Unlock()

		now := time.Now()
		data, exists := clients[key]

		if !exists || now.Sub(data.windowStart) > rateWindow {
			clients[key] = &clientData{count: 1, windowStart: now}
			c.Next()
			return
		}

		data.count++
		if data.count > rateLimit {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded. Try again later.",
			})
			return
		}

		c.Next()
	}
}

func resolveClientKey(c *gin.Context) string {
	tokenStr := c.GetHeader("Authorization")
	if tokenStr != "" {
		tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})
		if err == nil && token.Valid {
			if claims, ok := token.Claims.(jwt.MapClaims); ok {
				if userID, ok := claims["user_id"].(string); ok && userID != "" {
					return "user:" + userID
				}
			}
		}
	}
	return "ip:" + c.ClientIP()
}
