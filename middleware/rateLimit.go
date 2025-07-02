package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nelsonin-research-org/cdc-auth/globals"
	"github.com/nelsonin-research-org/cdc-auth/message"
)

func RateLimitMiddleware(limit int) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP() 
		endpoint := c.Request.URL.Path 
		
		globals.RequestStore.Mu.Lock()
		defer globals.RequestStore.Mu.Unlock()

		if _, exists := globals.RequestStore.Requests[ip]; !exists {
			globals.RequestStore.Requests[ip] = make(map[string]int)
		}

		if globals.RequestStore.Requests[ip][endpoint] >= limit {
			c.JSON(http.StatusTooManyRequests, message.ReturnMessage(http.StatusTooManyRequests))
			c.Abort() 
			return
		}

		globals.RequestStore.Requests[ip][endpoint]++

		go func(ip, endpoint string) {
			time.Sleep(1 * time.Hour)
			globals.RequestStore.Mu.Lock()
			defer globals.RequestStore.Mu.Unlock()
			globals.RequestStore.Requests[ip][endpoint] = 0
		}(ip, endpoint)

		c.Next() 
	}
}