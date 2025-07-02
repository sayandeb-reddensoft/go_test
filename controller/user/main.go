package controller

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/nelsonin-research-org/cdc-auth/interfaces"
	"github.com/nelsonin-research-org/cdc-auth/models/appschema"
	dbModel "github.com/nelsonin-research-org/cdc-auth/models/database"
	model "github.com/nelsonin-research-org/cdc-auth/models/user"
	"github.com/nelsonin-research-org/cdc-auth/utils"

	"gorm.io/gorm"
)

type userControllerImpl struct {
	RelationalDB *gorm.DB
}

func NewUserController(rdb *gorm.DB) interfaces.UserController {
	return &userControllerImpl{RelationalDB: rdb}
}

// user operations

func (s *userControllerImpl) CreateNewOrg(userData *model.SignUpData, userId string) (string, error) {
	addressInt, err := s.CreateNewAddress(&userData.Address)
	if err != nil {
		return "", err
	}

	newUserProfile := dbModel.Organization{
		OrgName: userData.OrgName,
		OrgUserID: userId,
		AddressID: addressInt,
		ContactName: userData.ContactName,
		ContactDescription: userData.ContactDescription,
		ContactNumber: userData.ContactNumber,
	}

	if err := s.RelationalDB.Create(&newUserProfile).Error; err != nil {
		log.Println("error while create the user :", err.Error())
		return "", errors.New("error while create new user")
	}

	return userId, nil
}

func (s *userControllerImpl) CreateNewAddress(address *model.SignupAddress) (uint, error) {
	zipCode, err := utils.StringToInt(address.PostalCode)
	if err != nil {
		return 0, err
	}

	newAddress := &dbModel.Address{
		StreetLine: address.StreetLane,
		City: address.City,
		State: address.State,
		PostalIndexCode: int32(zipCode),
	}

	affected := s.RelationalDB.Create(newAddress).RowsAffected
	if affected > 0 {
		return newAddress.ID, nil
	}

	return 0, nil
}

func (s *userControllerImpl) CreateNewUser(userData *model.User, userType int) (bool, error)  {
	hashedPass := utils.HashPassword(userData.Password)
	if hashedPass == "" {
		return false, errors.New("please provide a valid password")
	}

	newUser := dbModel.User{
		UserId: userData.UserId,
		Email: userData.Email,
		Password: hashedPass,
		RoleID: uint(userType),
		IsActive: false,
	}

	if err := s.RelationalDB.Create(&newUser).Error; err != nil {
		return false, err
	}

	return true, nil
}

func (s *userControllerImpl) IsUserAlreadyExists(userEmail string) (bool, error) {
	var userCount int64

	err := s.RelationalDB.Model(&dbModel.User{}).Where(dbModel.USER_COLLECTION_FIELDS.EMAIL+" = ?", utils.FormatStringToLowerCase(userEmail)).Count(&userCount).Error
	if err != nil {
		return false, err
	}

	if userCount > 0 {
		return true, nil
	}

	return false, nil
}

func (s *userControllerImpl) GetUserIdAndPasswordByEmail(email string) (string, string) {
	var user dbModel.User

	s.RelationalDB.Model(&dbModel.User{}).Select(dbModel.USER_COLLECTION_FIELDS.USER_ID, dbModel.USER_COLLECTION_FIELDS.PASSWORD).Where(dbModel.USER_COLLECTION_FIELDS.EMAIL, utils.FormatStringToLowerCase(email)).First(&user)

	return user.UserId, user.Password
}

func (s *userControllerImpl) LogThisLogin(userId string, isCredentialCorrect bool, userIp string) error {
	loginLog := &dbModel.Login{}

	// Fetch the existing record
	result := s.RelationalDB.Where("user_id = ?", userId).First(loginLog)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			loginLog = &dbModel.Login{
				UserId:       userId,
				LoginIp:      userIp,
				NoOfAttempts: 1,
			}

			if isCredentialCorrect {
				loginLog.LoginTimeStamp = time.Now()
				loginLog.LogOutTimeStamp = time.Time{}
			}

			result = s.RelationalDB.Create(loginLog)
			if result.Error != nil {
				return result.Error
			}
		} else {
			return result.Error
		}
	} else {
		if isCredentialCorrect {
			loginLog.LoginTimeStamp = time.Now()
			loginLog.LogOutTimeStamp = time.Time{}
		}
		loginLog.LoginIp = userIp
		loginLog.NoOfAttempts++

		result = s.RelationalDB.Model(&dbModel.Login{}).Where(dbModel.LOGIN_COLLECTION_FIELDS.USER_ID + "= ?", userId).Updates(loginLog)
		if result.Error != nil {
			return result.Error
		}
	}

	return nil
}

func (s *userControllerImpl) LogThisLogout(userId string) error {
	var count int64
	if err := s.RelationalDB.Model(&dbModel.Login{}).Where(dbModel.LOGIN_COLLECTION_FIELDS.USER_ID + " = ?", userId).Count(&count).Error; err != nil {
		return err
	}

	if count == 0 {
		return nil
	}

	result := s.RelationalDB.Model(&dbModel.Login{}).Where(dbModel.LOGIN_COLLECTION_FIELDS.USER_ID  + "= ?", userId).Update(dbModel.LOGIN_COLLECTION_FIELDS.LOGOUT_TIMESTAMP , time.Now())
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (s *userControllerImpl) GetUserDataFromSession(userId string) (*appschema.JwtData, error) {
	var user dbModel.User

	s.RelationalDB.Model(&dbModel.User{}).Where(dbModel.USER_COLLECTION_FIELDS.USER_ID, userId).First(&user)

	jwtData := appschema.JwtData{
		ID: userId,
		Email: user.Email,
		Role: int(user.RoleID),
	}

	return &jwtData, nil
}

func (s *userControllerImpl) GetUserDataById(userId string) (*model.GetUserData, error) {
	var user dbModel.User

	err := s.RelationalDB.Model(&dbModel.User{}).Where(dbModel.USER_COLLECTION_FIELDS.USER_ID + "= ?", userId).First(&user).Error
	if err != nil {
		return nil, errors.New("error fetching user data: " + err.Error())
	}

	userData := &model.GetUserData{
		Email: user.Email,
		Role: int(user.RoleID),
		IsActive: user.IsActive,
	}

	return userData, nil
}

func (s *userControllerImpl) UpdateUserPassword(password string, email string) (bool, error) {
	var user *dbModel.User

	hashedPass := utils.HashPassword(password)
	if hashedPass == "" {
		return false, fmt.Errorf("please choose different password")
	}

	res := s.RelationalDB.Model(&user).Where(dbModel.USER_COLLECTION_FIELDS.EMAIL + "= ?", utils.FormatStringToLowerCase(email)).Update(dbModel.USER_COLLECTION_FIELDS.PASSWORD, hashedPass)
	if res.RowsAffected == 0 {
		return false, fmt.Errorf("couldn't update the password")
	}

	return true, nil
}

func (s *userControllerImpl) IsUserAccountVerifiedByEmail(email string) (bool, error) {
	var count int64

	err := s.RelationalDB.Model(&dbModel.User{}).Where(dbModel.USER_COLLECTION_FIELDS.EMAIL + "= ? AND " + dbModel.USER_COLLECTION_FIELDS.IS_ACTIVE + " = ?", utils.FormatStringToLowerCase(email), true).Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("failed to check if user account is verified: %w", err)
	}

	return count > 0, nil
}

func (s *userControllerImpl) MakeUserAccountAsVerifiedByEmail(email string) (bool, error) {
	err := s.RelationalDB.Model(&dbModel.User{}).Where(dbModel.USER_COLLECTION_FIELDS.EMAIL + " = ?", email).Update(dbModel.USER_COLLECTION_FIELDS.IS_ACTIVE, true).Error
	if err != nil {
		return false, err
	}

	return true, nil
}
