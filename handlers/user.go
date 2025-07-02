package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	constants "github.com/nelsonin-research-org/cdc-auth/const"
	"github.com/nelsonin-research-org/cdc-auth/interfaces"
	"github.com/nelsonin-research-org/cdc-auth/message"
	"github.com/nelsonin-research-org/cdc-auth/models/appschema"
	emailModel "github.com/nelsonin-research-org/cdc-auth/models/email"
	limitations "github.com/nelsonin-research-org/cdc-auth/models/limitation"
	model "github.com/nelsonin-research-org/cdc-auth/models/user"
	userModel "github.com/nelsonin-research-org/cdc-auth/models/user"
	"github.com/nelsonin-research-org/cdc-auth/utils"
)

type UserHandler struct {
	UserService         interfaces.UserService
	FormValidateService interfaces.FormValidateService
	EmailService        interfaces.EmailService
	OTPService          interfaces.OTPService
}

func NewUserHandler(userService interfaces.UserService, formValidateService interfaces.FormValidateService, emailService interfaces.EmailService, otpService interfaces.OTPService) *UserHandler {
	return &UserHandler{
		UserService:         userService,
		FormValidateService: formValidateService,
		EmailService: 		 emailService,
		OTPService: 		 otpService,		
	}
}

func (h *UserHandler) CreateOrganization(c *gin.Context) {
	var signUpData userModel.SignUpData

	if err := c.ShouldBindJSON(&signUpData); err != nil {
		c.JSON(http.StatusBadRequest, message.ReturnInvalidFieldMsg())
		return
	}

	if err := h.FormValidateService.ValidateStruct(&signUpData); err != nil {
		invalidField := h.FormValidateService.ReturnFirstInvalidField(err)
		c.JSON(http.StatusBadRequest, message.ReturnCustomMessage(fmt.Sprintf("Invalid field: %s", invalidField)))
		return
	}

	exists, err := h.UserService.IsUserAlreadyExists(signUpData.OrgEmail)
	if err != nil {
		fmt.Println("error check while user is already exist :", err.Error())
		c.JSON(http.StatusInternalServerError, message.ReturnMessage(http.StatusInternalServerError))
		return
	} else if exists {
		c.JSON(http.StatusConflict, message.ReturnCustomMessage("email address already in use"))
		return
	}

	if !utils.IsStrongPassword(signUpData.Password) {
		c.JSON(http.StatusBadRequest, message.ReturnCustomMessage("Password must be 8 characters long, must include minimum one uppercase, one lowercase, one number and one special character"))
		return
	}

	userId := utils.GenerateUserUID()
	if userId == "" {
		fmt.Println("failed to generate user UID")
		c.JSON(http.StatusInternalServerError, message.ReturnMessage(http.StatusInternalServerError))
		return
	}

	userData := &model.User{
		UserId:   userId,
		Password: signUpData.Password,
		Email:    signUpData.OrgEmail,
	}

	ok, err := h.UserService.CreateNewUser(userData, int(constants.UserRole(constants.OrgAdmin)))
	if err != nil || !ok {
		fmt.Println("error creating user", err.Error())
		c.JSON(http.StatusInternalServerError, message.ReturnMessage(http.StatusInternalServerError))
		return
	}
	
	_, err = h.UserService.CreateNewOrg(&signUpData, userId)
	if err != nil {
		fmt.Println("error creating user :", err.Error())
		c.JSON(http.StatusInternalServerError, message.ReturnCustomMessage(err.Error()))
		return
	}

	otp, err := h.OTPService.GenerateOTP(limitations.OTP_LIMITATION.OTP_LENGTH)
	if err != nil {
		fmt.Println("error generating otp", err.Error())
		c.JSON(http.StatusInternalServerError, message.ReturnMessage(http.StatusInternalServerError))
		return 
	}

	ok, err = h.OTPService.StoreOtp(c.Request.Context(), limitations.OTP_LIMITATION.OTP_EXP, constants.MEMORY_VERIFY_KEY, otp)
	if err != nil || !ok {
		fmt.Println("error storing otp :", err.Error())
		c.JSON(http.StatusInternalServerError, message.ReturnMessage(http.StatusInternalServerError))
		return
	}

	if utils.IsStage() || utils.IsProduction() {
		mailContent := &emailModel.WelcomeAccountMailContent{
			Code: otp,
			Name: signUpData.OrgName,
		}

		ok, err := h.EmailService.SendWelcomeOTP(mailContent, signUpData.OrgEmail)
		if err != nil || !ok {
			fmt.Println("error sending otp:", err.Error())
			c.JSON(http.StatusInternalServerError, message.ReturnSomethingWentWrongMsg())
			return
		}
	} else {
		ok, err = h.OTPService.StoreOtp(c.Request.Context(), limitations.OTP_LIMITATION.OTP_EXP, (userId + constants.MEMORY_VERIFY_KEY), otp)
		if err != nil || !ok {
			fmt.Println("error storing otp :", err.Error())
			c.JSON(http.StatusInternalServerError, message.ReturnMessage(http.StatusInternalServerError))
			return
		}

		fmt.Println("verify otp : ", otp)
	}

	c.JSON(http.StatusOK, message.ReturnCustomDataWithKey("data", signUpData.OrgEmail))
}

func (h *UserHandler) Login(c *gin.Context) {
	var loginData userModel.LoginData

	if err := c.ShouldBindJSON(&loginData); err != nil {
		fmt.Println("error while check the login data :", err.Error())
		c.JSON(http.StatusBadRequest, message.ReturnCustomMessage("please provide a valid data"))
		return
	}

	if err := h.FormValidateService.ValidateStruct(&loginData); err != nil {
		msg := "invalid " + h.FormValidateService.ReturnFirstInvalidField(err)
		c.JSON(http.StatusBadRequest, message.ReturnCustomMessage(msg))
		return
	}
	
	userId, hashPwd := h.UserService.GetUserIdAndPasswordByEmail(loginData.Email)
	if userId == "" {
		c.JSON(http.StatusUnprocessableEntity, message.ReturnCustomMessage("email not found"))
		return
	}
	
	verified, err := h.UserService.IsUserAccountVerifiedByEmail(loginData.Email)
	if err != nil {
		fmt.Println("error checking user verification status:", err)
		c.JSON(http.StatusInternalServerError, message.ReturnSomethingWentWrongMsg())
		return
	}
	
	if !verified {
		c.JSON(http.StatusProxyAuthRequired, message.ReturnCustomMessage("email not verified."))
		return
	}

	userIp := c.ClientIP()

	isCorrect, errMsg := utils.VerifyPassword(hashPwd, loginData.Password)
	if !isCorrect {
		// log on failure attempt
		err := h.UserService.LogThisLogin(userId, false, userIp)
		if err != nil {
			fmt.Println("error log the login", err.Error())
			c.JSON(http.StatusInternalServerError, message.ReturnSomethingWentWrongMsg())
			return
		}
		c.JSON(http.StatusUnauthorized, message.ReturnCustomMessage(errMsg))
		return
	}

	// log on success attempt
	err = h.UserService.LogThisLogin(userId, true, userIp)
	if err != nil {
		fmt.Println("error log the login", err.Error())
		c.JSON(http.StatusInternalServerError, message.ReturnSomethingWentWrongMsg())
		return
	}

	jwtData, err := h.UserService.GetUserDataFromSession(userId)
	if err != nil {
		fmt.Println("error getting user data", err.Error())
		c.JSON(http.StatusInternalServerError, message.ReturnSomethingWentWrongMsg())
		return
	}

	tokens, err := utils.GeneratePrimaryToken(limitations.JWT_LIMITATION.PRIMARY_TOKEN_TTL, jwtData)
	if err != nil {
		fmt.Println("error creating tokens", err.Error())
		c.JSON(http.StatusInternalServerError, message.ReturnMessage(http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusOK, message.ReturnCustomDataWithKey("data", tokens))
}

func (h *UserHandler) Logout(c *gin.Context) {
	userId, err := utils.GetUserIdFromHeader(c)
	if err != nil || userId == "" {
		c.JSON(http.StatusUnauthorized, message.ReturnMessage(http.StatusUnauthorized))
		return
	}

	err = h.UserService.LogThisLogout(userId)
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusInternalServerError, message.ReturnMessage(http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusOK, message.ReturnMessage(http.StatusOK))
}

func (h *UserHandler) VerifyOtp(c *gin.Context) {
	var newOtpData userModel.VerifyOtpData

	if err := c.ShouldBindJSON(&newOtpData); err != nil {
		c.JSON(http.StatusBadRequest, message.ReturnInvalidFieldMsg())
		return
	}

	if err := h.FormValidateService.ValidateStruct(newOtpData); err != nil {
		msg := "Invalid field: " + h.FormValidateService.ReturnFirstInvalidField(err)
		c.JSON(http.StatusBadRequest, message.ReturnCustomMessage(msg))
		return
	}

	userId, _ := h.UserService.GetUserIdAndPasswordByEmail(newOtpData.Email)
	if userId == "" {
		c.JSON(http.StatusUnprocessableEntity, message.ReturnCustomMessage("email not found"))
		return
	}

	var otpKey string
	switch newOtpData.Type {
	case constants.RESET_REQUEST_TYPE:
		otpKey = userId + constants.MEMORY_FORGET_KEY
	case constants.VERIFY_REQUEST_TYPE:
		otpKey = userId + constants.MEMORY_VERIFY_KEY
	default:
		c.JSON(http.StatusBadRequest, message.ReturnCustomMessage("invalid OTP type"))
		return
	}

	verified, err := h.UserService.IsUserAccountVerifiedByEmail(newOtpData.Email)
	if err != nil {
		fmt.Println("error checking user verification status:", err)
		c.JSON(http.StatusInternalServerError, message.ReturnSomethingWentWrongMsg())
		return
	}
	
	switch newOtpData.Type {
	case constants.VERIFY_REQUEST_TYPE:
		if verified {
			c.JSON(http.StatusConflict, message.ReturnCustomMessage("email already verified."))
			return
		} else {
			if ok, err := h.OTPService.VerifyOTP(c.Request.Context(), otpKey, newOtpData.OTP); err != nil || !ok {
				c.JSON(http.StatusBadRequest, message.ReturnCustomMessage(err.Error()))
				return
			}

			ok, err := h.UserService.MakeUserAccountAsVerifiedByEmail(newOtpData.Email)
			if err != nil || !ok {
				fmt.Println("error making user account as verified :", err)
				c.JSON(http.StatusInternalServerError, message.ReturnMessage(http.StatusInternalServerError))
				return
			}
		}

		jwtData, err := h.UserService.GetUserDataFromSession(userId)
		if err != nil {
			fmt.Println("error creating token after otp verification :", err.Error())
			c.JSON(http.StatusInternalServerError, message.ReturnMessage(http.StatusInternalServerError))
			return
		}

		jwtToken, err := utils.GeneratePrimaryToken(limitations.JWT_LIMITATION.PRIMARY_TOKEN_TTL, jwtData)
		if err != nil {
			c.JSON(http.StatusInternalServerError, message.ReturnCustomMessage("failed to generate JWT token"))
			return
		}

		c.JSON(http.StatusOK, message.ReturnCustomDataWithKey("data", jwtToken))
	case constants.RESET_REQUEST_TYPE:
		if ok, err := h.OTPService.VerifyOTP(c.Request.Context(), otpKey, newOtpData.OTP); err != nil || !ok {
			c.JSON(http.StatusBadRequest, message.ReturnCustomMessage(err.Error()))
			return
		}

		if !verified {
			ok, err := h.UserService.MakeUserAccountAsVerifiedByEmail(newOtpData.Email)
			if err != nil || !ok {
				fmt.Println("error making user account as verified :", err)
				c.JSON(http.StatusInternalServerError, message.ReturnMessage(http.StatusInternalServerError))
				return
			}
		}

		data := &appschema.TempJwtData{
			Email: newOtpData.Email,
		}

		token, err := utils.GenerateTempToken(limitations.JWT_LIMITATION.TEMP_TOKEN_TTL, data)
		if err != nil {
			c.JSON(http.StatusInternalServerError, message.ReturnCustomMessage("failed to generate temp token"))
			return
		}

		c.JSON(http.StatusOK, message.ReturnCustomDataWithKey("data", token))
	default:
		c.JSON(http.StatusBadRequest, message.ReturnCustomMessage("invalid type"))
	}
}

func (h *UserHandler) ResendOTP(c *gin.Context) {
	var resendOtpData userModel.ResendOtpData

	if err := c.ShouldBindJSON(&resendOtpData); err != nil {
		c.JSON(http.StatusBadRequest, message.ReturnInvalidFieldMsg())
		return
	}

	if err := h.FormValidateService.ValidateStruct(resendOtpData); err != nil {
		msg := "Invalid field: " + h.FormValidateService.ReturnFirstInvalidField(err)
		c.JSON(http.StatusBadRequest, message.ReturnCustomMessage(msg))
		return
	}

	userId, _ := h.UserService.GetUserIdAndPasswordByEmail(resendOtpData.Email)
	if userId == "" {
		c.JSON(http.StatusUnprocessableEntity, message.ReturnCustomMessage("email not found"))
		return
	}
	
	verified, err := h.UserService.IsUserAccountVerifiedByEmail(resendOtpData.Email)
	if err != nil {
		fmt.Println("error checking user verification status:", err)
		c.JSON(http.StatusInternalServerError, message.ReturnSomethingWentWrongMsg())
		return
	}

	var otpKey string
	switch resendOtpData.Type {
	case constants.RESET_REQUEST_TYPE:
		otpKey = userId + constants.MEMORY_FORGET_KEY
	case constants.VERIFY_REQUEST_TYPE:
		otpKey = userId + constants.MEMORY_VERIFY_KEY
	default:
		c.JSON(http.StatusBadRequest, message.ReturnCustomMessage("invalid OTP type"))
		return
	}

	otp, err := h.OTPService.GenerateOTP(limitations.OTP_LIMITATION.OTP_LENGTH)
	if err != nil {
		fmt.Println("error generating otp", err.Error())
		c.JSON(http.StatusInternalServerError, message.ReturnMessage(http.StatusInternalServerError))
		return 
	}

	if utils.IsStage() || utils.IsProduction() {
		switch resendOtpData.Type {
		case constants.RESET_REQUEST_TYPE:
			ok, err := h.OTPService.StoreOtp(c.Request.Context(), limitations.OTP_LIMITATION.OTP_EXP, otpKey, otp)
			if err != nil || !ok {
				fmt.Println("error storing otp :", err.Error())
				c.JSON(http.StatusInternalServerError, message.ReturnMessage(http.StatusInternalServerError))
				return
			}

			mailContent := &emailModel.ResetPasswordMailContent{
				Code: otp,
			}

			ok, err = h.EmailService.SendResetPasswordOTP(mailContent, resendOtpData.Email)
			if err != nil || !ok {
				fmt.Println("filed to resend otp ", err.Error())
				c.JSON(http.StatusInternalServerError, message.ReturnSomethingWentWrongMsg())
				return
			}
		case constants.VERIFY_REQUEST_TYPE:
			if verified {
				c.JSON(http.StatusConflict, message.ReturnCustomMessage("email is already verified."))
				return
			}

			ok, err := h.OTPService.StoreOtp(c.Request.Context(), limitations.OTP_LIMITATION.OTP_EXP, otpKey, otp)
			if err != nil || !ok {
				fmt.Println("error storing otp :", err.Error())
				c.JSON(http.StatusInternalServerError, message.ReturnMessage(http.StatusInternalServerError))
				return
			}

			user, err := h.UserService.GetUserDataById(userId)
			if err != nil {
				c.JSON(http.StatusInternalServerError, message.ReturnSomethingWentWrongMsg())
				return
			}
			
			mailContent := &emailModel.WelcomeAccountMailContent{
				Code: otp,
				Name: user.Name,
			}
			
			ok, err = h.EmailService.SendWelcomeOTP(mailContent, resendOtpData.Email)
			if err != nil || !ok {
				fmt.Println("filed to resend otp ", err.Error())
				c.JSON(http.StatusInternalServerError, message.ReturnSomethingWentWrongMsg())
				return
			}
		default:
			c.JSON(http.StatusBadRequest, message.ReturnCustomMessage("invalid OTP type"))
			return
		}
	} else {
		if verified && resendOtpData.Type == constants.VERIFY_REQUEST_TYPE {
			c.JSON(http.StatusConflict, message.ReturnCustomMessage("email is already verified."))
			return
		}

		ok, err := h.OTPService.StoreOtp(c.Request.Context(), limitations.OTP_LIMITATION.OTP_EXP, otpKey, otp)
		if err != nil || !ok {
			fmt.Println("error storing otp :", err.Error())
			c.JSON(http.StatusInternalServerError, message.ReturnMessage(http.StatusInternalServerError))
			return
		}

		fmt.Println("resend otp : ", otp)
	}

	c.JSON(http.StatusOK, message.ReturnMessage(http.StatusOK))
}

func (h *UserHandler) RefreshTokenHandler(c *gin.Context) {
	userId, err := utils.GetUserIdFromHeader(c)
	if err != nil || userId == "" {
		c.JSON(http.StatusUnauthorized, message.ReturnMessage(http.StatusUnauthorized))
		return
	}

	jwtData, err := h.UserService.GetUserDataFromSession(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, message.ReturnSomethingWentWrongMsg())
		return
	}

	tokens, err := utils.GeneratePrimaryToken(limitations.JWT_LIMITATION.PRIMARY_TOKEN_TTL, jwtData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, message.ReturnCustomMessage("failed to generte token"))
		return
	}

	c.JSON(http.StatusOK, message.ReturnCustomDataWithKey("data", tokens))
}

func (h *UserHandler) ForgetPassword(c *gin.Context) {
	var fpData userModel.ForgetPasswordData

	if err := c.ShouldBindJSON(&fpData); err != nil {
		c.JSON(http.StatusBadRequest, message.ReturnInvalidFieldMsg())
		return
	}

	if err := h.FormValidateService.ValidateStruct(fpData); err != nil {
		msg := "invalid " + h.FormValidateService.ReturnFirstInvalidField(err)
		c.JSON(http.StatusBadRequest, message.ReturnCustomMessage(msg))
		return
	}

	userId, _ := h.UserService.GetUserIdAndPasswordByEmail(fpData.Email)
	if userId == "" {
		c.JSON(http.StatusUnprocessableEntity, message.ReturnCustomMessage("email not found"))
		return
	}
	
	otp, err := h.OTPService.GenerateOTP(limitations.OTP_LIMITATION.OTP_LENGTH)
	if err != nil {
		fmt.Println("error generating otp", err.Error())
		c.JSON(http.StatusInternalServerError, message.ReturnMessage(http.StatusInternalServerError))
		return 
	}

	if utils.IsStage() || utils.IsProduction() {
		mailContent := &emailModel.ResetPasswordMailContent{
			Code: otp,
		}

		ok, err := h.EmailService.SendResetPasswordOTP(mailContent, fpData.Email)
		if err != nil || !ok {
			fmt.Println("filed to send otp ", err.Error())
			c.JSON(http.StatusInternalServerError, message.ReturnSomethingWentWrongMsg())
			return
		}
	} else {
		ok, err := h.OTPService.StoreOtp(c.Request.Context(), limitations.OTP_LIMITATION.OTP_EXP, (userId+constants.MEMORY_FORGET_KEY), otp)
		if err != nil || !ok {
			fmt.Println("error storing otp :", err.Error())
			c.JSON(http.StatusInternalServerError, message.ReturnMessage(http.StatusInternalServerError))
			return
		}

		fmt.Println("reset otp : ", otp)
	}

	c.JSON(http.StatusOK, message.ReturnMessage(http.StatusOK))
}

func (h *UserHandler) UpdatePassword(c *gin.Context) {
	var updatePassData userModel.UpdatePasswordData

	email, err := utils.GetUserEmailFromHeader(c)
	if err != nil || email == "" {
		c.JSON(http.StatusUnauthorized, message.ReturnMessage(http.StatusUnauthorized))
		return
	}
	
	if err := c.ShouldBindJSON(&updatePassData); err != nil {
		c.JSON(http.StatusBadRequest, message.ReturnInvalidFieldMsg())
		return
	}

	if err := h.FormValidateService.ValidateStruct(updatePassData); err != nil {
		msg := "invalid " + h.FormValidateService.ReturnFirstInvalidField(err)
		c.JSON(http.StatusBadRequest, message.ReturnCustomMessage(msg))
		return
	}

	if !utils.IsStrongPassword(updatePassData.Password) {
		c.JSON(http.StatusBadRequest, message.ReturnCustomMessage("New Password must be 8 characters long, must include minimum one uppercase, one lowercase, one number and one special character"))
		return
	}

	_, hashPwd := h.UserService.GetUserIdAndPasswordByEmail(email)
	if hashPwd == "" {
		c.JSON(http.StatusUnprocessableEntity, message.ReturnCustomMessage("email not found"))
		return
	}

	if utils.MatchWithHashPassword(updatePassData.Password, hashPwd) {
		c.JSON(http.StatusBadRequest, message.ReturnCustomMessage("new password should not be same as old password"))
		return
	}

	ok, err := h.UserService.UpdateUserPassword(updatePassData.Password, email)
	if err != nil || !ok {
		c.JSON(http.StatusInternalServerError, message.ReturnCustomMessage("unable to update the password"))
		return
	}

	c.JSON(http.StatusOK, message.ReturnMessage(http.StatusOK))
}

