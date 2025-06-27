package handlers

import (
	"github.com/gin-gonic/gin"

	controller "github.com/nelsonin-research-org/clenz-auth/controller/user"
	"github.com/nelsonin-research-org/clenz-auth/globals"
	"github.com/nelsonin-research-org/clenz-auth/services"
)

type UserHandler struct {
	UserService         services.UserServiceImpl
	FormValidateService services.FormValidationServiceImpl
	EmailService        services.EmailServiceImpl
}

func NewUserHandler() *UserHandler {
	userController := &controller.UserController{
		RelationalDB: globals.RelationalDb,
		RedisDB:      globals.RedisClient,
	}

	return &UserHandler{
		UserService: services.UserServiceImpl{
			Controller: *userController,
		},
		FormValidateService: services.FormValidationServiceImpl{},
		EmailService:        services.EmailServiceImpl{},
	}
}

func (h *UserHandler) SignUp(c *gin.Context) {
	 
}

func (h *UserHandler) Login(c *gin.Context) {
	
}

func (h *UserHandler) Logout(c *gin.Context) {

}

func (h *UserHandler) VerifyOtp(c *gin.Context) {

}

func (h *UserHandler) ResendOTP(c *gin.Context) {

}

func (h *UserHandler) RefreshTokenHandler(c *gin.Context) {

}

func (h *UserHandler) ForgetPassword(c *gin.Context) {
	
}

func (h *UserHandler) UpdatePassword(c *gin.Context) {
	
}

