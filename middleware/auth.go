package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	constants "github.com/nelsonin-research-org/cdc-auth/const"
	"github.com/nelsonin-research-org/cdc-auth/globals"
	"github.com/nelsonin-research-org/cdc-auth/interfaces"
	"github.com/nelsonin-research-org/cdc-auth/message"
	"github.com/nelsonin-research-org/cdc-auth/utils"
)

func AuthMiddleware(userService interfaces.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, message.ReturnMessage(http.StatusUnauthorized))
			c.Abort()
			return
		}

		parts := strings.Fields(authHeader)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusBadRequest, message.ReturnMessage(http.StatusBadRequest))
			c.Abort()
			return
		}

		tokenString := parts[1]

		if globals.AppKeys.PrivateKey == nil || globals.AppKeys.PublicKeyPem == nil {
			fmt.Println("pem keys not loaded")
			c.JSON(http.StatusInternalServerError, message.ReturnMessage(http.StatusInternalServerError))
			c.Abort()
			return
		}

		claims, err := utils.ParseToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, message.ReturnCustomMessage("invalid token"))
			c.Abort()
			return
		}
		
		tokenTypeRaw, ok := claims["token_type"]
		if !ok {
			c.JSON(http.StatusBadRequest, message.ReturnCustomMessage("token type not found in token"))
			c.Abort()
			return
		}
		
		tokenTypeFloat, ok := tokenTypeRaw.(float64)
		if !ok {
			c.JSON(http.StatusBadRequest, message.ReturnCustomMessage("invalid token type format"))
			c.Abort()
			return
		}
		
		if int(tokenTypeFloat) != constants.PRIMARY_TOKEN && int(tokenTypeFloat) != constants.REFRESH_TOKEN {
			c.JSON(http.StatusBadRequest, message.ReturnCustomMessage("token type not matched"))
			c.Abort()
			return
		}

		userId, _ := userService.GetUserIdAndPasswordByEmail(claims["email"].(string))
		if userId == "" {
			c.JSON(http.StatusUnprocessableEntity, message.ReturnCustomMessage("acount not found"))
			return
		}
				
		verified, err := userService.IsUserAccountVerifiedByEmail(claims["email"].(string))
		if err != nil {
			c.JSON(http.StatusInternalServerError, message.ReturnMessage(http.StatusInternalServerError))
			c.Abort()
			return
		}
		
		if !verified {
			c.JSON(http.StatusProxyAuthRequired, message.ReturnCustomMessage("acount not verified"))
			c.Abort()
			return
		}

		c.Set("id", claims["id"])
		c.Set("email", claims["email"])
		c.Set("role", claims["role"])
		c.Next()
	}
}