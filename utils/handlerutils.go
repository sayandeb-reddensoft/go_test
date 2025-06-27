package utils

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func GetUserIdFromHeader(c *gin.Context) (string, error) {
	userId, exists := c.Get("id")
	if !exists {
		return "", fmt.Errorf("userId is not exist")
	}

	return userId.(string), nil
}

func GetUserEmailFromHeader(c *gin.Context) (string, error) {
	email, exists := c.Get("email")
	if !exists {
		return "", fmt.Errorf("user email is not exist")
	}

	return email.(string), nil
}
