package routes

import (
	"github.com/gin-gonic/gin"
	app "github.com/nelsonin-research-org/cdc-auth/handlers/data"
	"github.com/nelsonin-research-org/cdc-auth/middleware"
	limitations "github.com/nelsonin-research-org/cdc-auth/models/limitation"
)

func NoAuthGroupRoutes(r *gin.RouterGroup, handlers *app.AppHandlers) {
	// user
	r.POST("/login", handlers.UserHandler.Login)
	r.POST("/org/create", handlers.UserHandler.CreateOrganization)
	r.POST("/verify/OTP", handlers.UserHandler.VerifyOtp)
	r.POST("/resend/OTP", middleware.RateLimitMiddleware(limitations.HANDLER_LIMITATION.RESEND_OTP), handlers.UserHandler.ResendOTP)
	r.POST("/forget-password", middleware.RateLimitMiddleware(limitations.HANDLER_LIMITATION.FORGOT_PASSWORD), handlers.UserHandler.ForgetPassword)
	r.PUT("/password", middleware.ValidateTempToken(handlers.UserHandler.UserService), middleware.RateLimitMiddleware(limitations.HANDLER_LIMITATION.UPDATE_PASSWORD), handlers.UserHandler.UpdatePassword)
	r.POST("/account/logout", middleware.AuthMiddleware(handlers.UserHandler.UserService), handlers.UserHandler.Logout)

	// token
	r.GET("/refresh-token", middleware.AuthMiddleware(handlers.UserHandler.UserService), handlers.UserHandler.RefreshTokenHandler)
}
