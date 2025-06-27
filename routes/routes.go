package routes

import (
	"github.com/gin-gonic/gin"
	app "github.com/nelsonin-research-org/clenz-auth/handlers/data"
	"github.com/nelsonin-research-org/clenz-auth/middleware"
)

func NoAuthGroupRoutes(r *gin.RouterGroup, handlers app.AppHandlers) {
	// user
	r.POST("/login", handlers.UserHandler.Login)
	r.POST("/sign-up", handlers.UserHandler.SignUp)
	r.POST("/verify/OTP", handlers.UserHandler.VerifyOtp)
	r.POST("/resend/OTP", middleware.RateLimitMiddleware(5), handlers.UserHandler.ResendOTP)
	r.POST("/forget-password", middleware.RateLimitMiddleware(10), handlers.UserHandler.ForgetPassword)
	r.PUT("/password", middleware.ValidateTempToken(&handlers.UserHandler.UserService), middleware.RateLimitMiddleware(5), handlers.UserHandler.UpdatePassword)
	r.POST("/account/logout", middleware.AuthMiddleware(&handlers.UserHandler.UserService), handlers.UserHandler.Logout)

	// token
	r.GET("/refresh-token", middleware.AuthMiddleware(&handlers.UserHandler.UserService), handlers.UserHandler.RefreshTokenHandler)
}
