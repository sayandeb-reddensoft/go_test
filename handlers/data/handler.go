package handlers

import (
	"os"

	emailController "github.com/nelsonin-research-org/cdc-auth/controller/email"
	formValidationController "github.com/nelsonin-research-org/cdc-auth/controller/form"
	otpController "github.com/nelsonin-research-org/cdc-auth/controller/otp"
	userController "github.com/nelsonin-research-org/cdc-auth/controller/user"
	"github.com/nelsonin-research-org/cdc-auth/globals"
	"github.com/nelsonin-research-org/cdc-auth/handlers"
	"github.com/nelsonin-research-org/cdc-auth/services"
)

type AppHandlers struct {
	UserHandler     *handlers.UserHandler
}

func LoadAppHandlers() *AppHandlers {
	// user
	userController := userController.NewUserController(globals.RelationalDb)
	userService := services.NewUserService(userController)

	// form validation
	formValidationController := formValidationController.NewFormValidationController()
	formValidationService := services.NewFormValidationService(formValidationController)

	// otp
	otpController := otpController.NewOTPController(globals.RedisClient)
	otpService := services.NewOTPService(otpController)

	// email
	emailController := emailController.NewEmailController(
		os.Getenv("MAIL_TEMPLATE_PATH"), 
		os.Getenv("MAIL_HOST"), 
		os.Getenv("MAIL_USERNAME"),
		os.Getenv("MAIL_PASSWORD"), 
		os.Getenv("MAIL_FROM"), 
		os.Getenv("MAIL_PORT"), 
		globals.AWSSesSession,
	)
	emailService := services.NewEmailService(emailController)

	return &AppHandlers{
		UserHandler:          handlers.NewUserHandler(userService, formValidationService, emailService, otpService),
	}
}