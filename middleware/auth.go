package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/nelsonin-research-org/clenz-auth/services"
)

func AuthMiddleware(userService *services.UserServiceImpl) gin.HandlerFunc {
	return func(c *gin.Context) {
		
		c.Next()
	}
}