package services

import (
	"github.com/nelsonin-research-org/cdc-auth/interfaces"
	"github.com/nelsonin-research-org/cdc-auth/models/appschema"
	userModel "github.com/nelsonin-research-org/cdc-auth/models/user"
)

type userControllerImpl struct {
	userController interfaces.UserController
}

func NewUserService(c interfaces.UserController) interfaces.UserService {
	return &userControllerImpl{userController: c}
}

func (service *userControllerImpl) CreateNewOrg(userData *userModel.SignUpData, userId string) (string, error) {
	return service.userController.CreateNewOrg(userData, userId)
}

func (service *userControllerImpl) CreateNewUser(userData *userModel.User, uType int) (bool, error)  {
	return service.userController.CreateNewUser(userData, uType)
}

func (service *userControllerImpl) IsUserAlreadyExists(email string) (bool, error) {
	return service.userController.IsUserAlreadyExists(email)
}

func (service *userControllerImpl) GetUserIdAndPasswordByEmail(email string) (string, string) {
	return service.userController.GetUserIdAndPasswordByEmail(email)
}

func (service *userControllerImpl) LogThisLogin(userId string, isCredentialCorrect bool, userIp string) error {
	return service.userController.LogThisLogin(userId, isCredentialCorrect, userIp)
}

func (service *userControllerImpl) LogThisLogout(userId string) error {
	return service.userController.LogThisLogout(userId)
}

func (service *userControllerImpl) GetUserDataFromSession(userId string) (*appschema.JwtData, error) {
	return service.userController.GetUserDataFromSession(userId)
}

func (service *userControllerImpl) GetUserDataById(userId string) (*userModel.GetUserData, error) {
	return service.userController.GetUserDataById(userId)
}

func (service *userControllerImpl) UpdateUserPassword(password string, email string) (bool, error) {
	return service.userController.UpdateUserPassword(password, email)
}

func (service *userControllerImpl) IsUserAccountVerifiedByEmail(email string) (bool, error) {
	return service.userController.IsUserAccountVerifiedByEmail(email)
}

func (service *userControllerImpl) MakeUserAccountAsVerifiedByEmail(email string) (bool, error) {
	return service.userController.MakeUserAccountAsVerifiedByEmail(email)
}
