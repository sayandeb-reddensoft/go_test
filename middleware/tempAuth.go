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

func ValidateTempToken(userService interfaces.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			fmt.Println("Authorization header not exist")
			c.JSON(http.StatusUnauthorized, message.ReturnMessage(http.StatusUnauthorized))
			c.Abort()
			return
		}

		parts := strings.Fields(authHeader)
		if len(parts) != 2 || parts[0] != "Bearer" {
			fmt.Println("toekn not exist on header or invalid header")
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

		if int(tokenTypeFloat) != constants.TEMP_TOKEN {
			c.JSON(http.StatusBadRequest, message.ReturnCustomMessage("token type not matched"))
			c.Abort()
			return
		}
		
		exists, err := userService.IsUserAlreadyExists(claims["email"].(string))
		if err != nil {
			fmt.Println("error verifying user already exist or not")
			c.JSON(http.StatusInternalServerError, message.ReturnMessage(http.StatusInternalServerError))
			c.Abort()
			return
		}

		if !exists {
			c.JSON(http.StatusUnprocessableEntity, message.ReturnCustomMessage("user not found"))
			c.Abort()
			return
		}

		c.Set("email", claims["email"])
		c.Next()
	}
}
