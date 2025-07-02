package interfaces

import (
	"context"
	"time"

	"github.com/nelsonin-research-org/cdc-auth/models/appschema"
	emailModel "github.com/nelsonin-research-org/cdc-auth/models/email"
	userModel "github.com/nelsonin-research-org/cdc-auth/models/user"
)

type UserController interface {
	LogThisLogout(userId string) error   
	IsUserAlreadyExists(email string) (bool, error)
	IsUserAccountVerifiedByEmail(email string) (bool, error) 
	GetUserIdAndPasswordByEmail(email string) (string, string)
	MakeUserAccountAsVerifiedByEmail(email string) (bool, error)
	GetUserDataById(userId string) (*userModel.GetUserData, error)
	UpdateUserPassword(password string, email string) (bool, error) 
	CreateNewUser(userData *userModel.User, uType int) (bool, error) 
	GetUserDataFromSession(userId string) (*appschema.JwtData, error)
	LogThisLogin(userId string, isCredentialCorrect bool, userIp string) error
	CreateNewOrg(userData *userModel.SignUpData, userId string) (string, error) 
}

type EmailController interface {
	SendWelcomeOTP(data *emailModel.WelcomeAccountMailContent, emails ...string) (bool, error)
	SendResetPasswordOTP(data *emailModel.ResetPasswordMailContent, emails ...string) (bool, error)
	SendDeleteAccountOTP(data *emailModel.DeleteAccountMailContent, emails ...string) (bool, error) 
}

type OTPController interface {
	GenerateOTP(length int) (int, error)
	VerifyOTP(ctx context.Context, otpKey string, otp int) (bool, error)
	StoreOtp(ctx context.Context, ttl time.Duration, key string, otp int) (bool, error)
}

type FormValidationController interface {
	ValidateStruct(s interface{}) error 
	ReturnFirstInvalidField(err error) string 
}
