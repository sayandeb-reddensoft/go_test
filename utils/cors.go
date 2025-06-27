package utils

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func GetCorsConfig() gin.HandlerFunc {
	var origins = []string{"*"}

	return cors.New(cors.Config{

		AllowOrigins: origins,
		AllowMethods: []string{"GET", "POST", "DELETE", "OPTIONS", "PUT"},
		AllowHeaders: []string{"Authorization", "Accept", "Accept-Encoding",
			"Accept-Language", "Connection", "Content-Length",
			"Content-Type", "Host", "Origin", "Referer", "User-Agent"},
		AllowCredentials: true,
		MaxAge:           24 * time.Hour,
	})
}
