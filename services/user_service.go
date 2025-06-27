package services

import (
	controller "github.com/nelsonin-research-org/clenz-auth/controller/user"
	"github.com/nelsonin-research-org/clenz-auth/models/appschema"
	socialModel "github.com/nelsonin-research-org/clenz-auth/models/socialAuth"
	userModel "github.com/nelsonin-research-org/clenz-auth/models/user"
)

type UserServiceImpl struct {
	Controller controller.UserController
}

func NewUserService() *UserServiceImpl {
	return &UserServiceImpl{}
}

func (service *UserServiceImpl) CreateNewUser(userData *userModel.SignUpData) (string, error) {
	return service.Controller.CreateNewUser(userData)
}

func (service *UserServiceImpl) IsUserAlreadyExists(email string) (bool, error) {
	return service.Controller.IsUserAlreadyExists(email)
}

func (service *UserServiceImpl) GetUserIdAndPasswordByEmail(email string) (string, string) {
	return service.Controller.GetUserIdAndPasswordByEmail(email)
}

func (service *UserServiceImpl) LogThisLogin(userId string, isCredentialCorrect bool, userIp string) error {
	return service.Controller.LogThisLogin(userId, isCredentialCorrect, userIp)
}

func (service *UserServiceImpl) LogThisLogout(userId string) error {
	return service.Controller.LogThisLogout(userId)
}

func (service *UserServiceImpl) GetUserDataFromSession(userId string) (*appschema.JwtData, error) {
	return service.Controller.GetUserDataFromSession(userId)
}

func (service *UserServiceImpl) GetUserDataById(userId string) (*userModel.GetUserData, error) {
	return service.Controller.GetUserDataById(userId)
}

func (service *UserServiceImpl) GenerateAndStoreOtp(userId string, key string) (int, error) {
	return service.Controller.GenerateAndStoreOtp(userId, key)
}

func (service *UserServiceImpl) VerifyOTP(otpKey string, otp int) (bool, error) {
	return service.Controller.VerifyOTP(otpKey, otp)
}

func (service *UserServiceImpl) UpdateUserPassword(password string, email string) (bool, error) {
	return service.Controller.UpdateUserPassword(password, email)
}

func (service *UserServiceImpl) IsUserAccountDeletedByEmail(email string) (bool, error) {
	return service.Controller.IsUserAccountDeletedByEmail(email)
}

func (service *UserServiceImpl) IsSocialUserExistAndActive(email string) (string, error) {
	return service.Controller.IsSocialUserExistAndActive(email)
}

func (service *UserServiceImpl) MakeUserAccountActiveByEmail(email string) (bool, error) {
	return service.Controller.MakeUserAccountActiveByEmail(email)
}

func (service *UserServiceImpl) IsUserAccountVerifiedByEmail(email string) (bool, error) {
	return service.Controller.IsUserAccountVerifiedByEmail(email)
}

func (service *UserServiceImpl) MakeUserAccountAsVerifiedByEmail(email string) (bool, error) {
	return service.Controller.MakeUserAccountAsVerifiedByEmail(email)
}

func (service *UserServiceImpl) CreateNewSocialUser(userData *socialModel.SocialSignupData) (string, error) {
	return service.Controller.CreateNewSocialUser(userData)
}

func (service *UserServiceImpl) SoftDeleteUserProfile(userId string, ch chan bool) {
	service.Controller.SoftDeleteUserProfile(userId, ch)
}

func (service *UserServiceImpl) UpsertDeleteRequet(reason, email, userId string) (bool, error) {
	return service.Controller.UpsertDeleteRequet(reason, email, userId)
}
