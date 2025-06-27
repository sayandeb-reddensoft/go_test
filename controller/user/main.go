package controller

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	constants "github.com/nelsonin-research-org/clenz-auth/const"
	"github.com/nelsonin-research-org/clenz-auth/globals"
	"github.com/nelsonin-research-org/clenz-auth/models/appschema"
	dbModel "github.com/nelsonin-research-org/clenz-auth/models/database"
	socialModel "github.com/nelsonin-research-org/clenz-auth/models/socialAuth"
	model "github.com/nelsonin-research-org/clenz-auth/models/user"
	"github.com/nelsonin-research-org/clenz-auth/utils"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type UserController struct {
	RelationalDB *gorm.DB
	RedisDB      *redis.Client
}

func NewUserController() *UserController {
	return &UserController{
		RelationalDB: globals.RelationalDb,
		RedisDB:      globals.RedisClient,
	}
}

var ctx = context.Background()

// user operations

func (s *UserController) CreateNewUser(userData *model.SignUpData) (string, error) {
	
	hashedPass := utils.HashPassword(userData.Password)
	if hashedPass == "" {
		return "", errors.New("please provide a valid password")
	}

	dob, err := utils.FormatStringToDate(userData.DOB)
	if err != nil {
		return "", errors.New("error while format date")
	}

	userUID := utils.GenerateUserUID()
	if userUID == "" {
		return "", errors.New("failed to generate user UID")
	}

	newUserProfile := dbModel.User{
		UserId:    userUID,
		FirstName: userData.FirstName,
		LastName:  userData.LastName,
		Password:  hashedPass,
		Type:      dbModel.UserType(constants.EMAIL_USER_TYPE),
		DOB:       dob,
		Email:     utils.FormatStringToLowerCase(userData.Email),
	}

	if err := s.RelationalDB.Create(&newUserProfile).Error; err != nil {
		log.Println("error while create the user :", err.Error())
		return "", errors.New("error while create new user")
	}

	return userUID, nil
}

func (s *UserController) IsUserAlreadyExists(userEmail string) (bool, error) {
	var userCount int64

	err := s.RelationalDB.Model(&dbModel.User{}).Where("email = ?", utils.FormatStringToLowerCase(userEmail)).Count(&userCount).Error
	if err != nil {
		return false, err
	}

	if userCount > 0 {
		return true, nil
	}

	return false, nil
}

func (s *UserController) GetUserIdAndPasswordByEmail(email string) (string, string) {
	var user dbModel.User

	s.RelationalDB.Model(&dbModel.User{}).Select("user_id", "password").Where("email", utils.FormatStringToLowerCase(email)).First(&user)

	return user.UserId, user.Password
}

func (s *UserController) LogThisLogin(userId string, isCredentialCorrect bool, userIp string) error {
	loginLog := &dbModel.Logins{}

	// Fetch the existing record
	result := s.RelationalDB.Where("user_id = ?", userId).First(loginLog)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			loginLog = &dbModel.Logins{
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

		result = s.RelationalDB.Model(&dbModel.Logins{}).Where("user_id = ?", userId).Updates(loginLog)
		if result.Error != nil {
			return result.Error
		}
	}

	return nil
}

func (s *UserController) LogThisLogout(userId string) error {
	var count int64
	if err := s.RelationalDB.Model(&dbModel.Logins{}).Where("user_id = ?", userId).Count(&count).Error; err != nil {
		return err
	}

	if count == 0 {
		return nil
	}

	result := s.RelationalDB.Model(&dbModel.Logins{}).Where("user_id = ?", userId).Update("log_out_time_stamp", time.Now())
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (s *UserController) GetUserDataFromSession(userId string) (*appschema.JwtData, error) {
	var jwtData appschema.JwtData
	var user dbModel.User

	s.RelationalDB.Model(&dbModel.User{}).Where("user_id", userId).First(&user)

	jwtData = appschema.JwtData{
		ID:        user.UserId,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Type:      int(user.Type),
	}

	return &jwtData, nil
}

func (s *UserController) GetUserDataById(userId string) (*model.GetUserData, error) {
	var user dbModel.User

	err := s.RelationalDB.Model(&dbModel.User{}).Where("user_id = ?", userId).First(&user).Error
	if err != nil {
		return nil, errors.New("error fetching user data: " + err.Error())
	}

	userData := model.GetUserData{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Type:      int(user.Type),
		DOB:       utils.FormatDateToString(user.DOB),
		IsActive:  user.IsActive,
		Gender:    user.Gender,
	}

	return &userData, nil
}

func (s *UserController) UpdateUserPassword(password string, email string) (bool, error) {
	var user *dbModel.User

	hashedPass := utils.HashPassword(password)
	if hashedPass == "" {
		return false, fmt.Errorf("please choose different password")
	}

	res := s.RelationalDB.Model(&user).Where("email = ?", utils.FormatStringToLowerCase(email)).Update("password", hashedPass)
	if res.RowsAffected == 0 {
		return false, fmt.Errorf("couldn't update the password")
	}

	return true, nil
}

func (s *UserController) IsUserAccountDeletedByEmail(email string) (bool, error) {
	var count int64

	err := s.RelationalDB.Model(&dbModel.User{}).Where("email = ? AND is_delete = ?", utils.FormatStringToLowerCase(email), true).Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("failed to check if user account is deleted: %w", err)
	}

	return count > 0, nil
}

func (s *UserController) IsUserAccountVerifiedByEmail(email string) (bool, error) {
	var count int64

	err := s.RelationalDB.Model(&dbModel.User{}).Where("email = ? AND is_active = ?", utils.FormatStringToLowerCase(email), true).Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("failed to check if user account is verified: %w", err)
	}

	return count > 0, nil
}

func (s *UserController) MakeUserAccountActiveByEmail(email string) (bool, error) {

	err := s.RelationalDB.Model(&dbModel.User{}).Where("email = ?", utils.FormatStringToLowerCase(email)).Update("is_delete", false).Error
	if err != nil {
		return false, err
	}
	
	return true, nil
}

func (s *UserController) MakeUserAccountAsVerifiedByEmail(email string) (bool, error) {
	err := s.RelationalDB.Model(&dbModel.User{}).Where("email = ?", email).Update("is_active", true).Error
	if err != nil {
		return false, err
	}

	return true, nil
}

// social login operations

func (s *UserController) CreateNewSocialUser(userData *socialModel.SocialSignupData) (string, error) {

	email := utils.FormatStringToLowerCase(userData.Email)
	if email == "" {
		return "", errors.New("email not exist")
	}

	userUID := utils.GenerateUserUID()
	if userUID == "" {
		return "", errors.New("failed to generate user UID")
	}

	dob, err := utils.FormatStringToDate(userData.DOB)
	if err != nil {
		return "", err
	}

	newUserProfile := dbModel.User{
		UserId:    userUID,
		FirstName: userData.FirstName,
		LastName:  userData.LastName,
		Email:     utils.FormatStringToLowerCase(userData.Email),
		DOB:       dob,
		IsActive:  true,
		Type:      dbModel.UserType(userData.Type),
	}

	if err := s.RelationalDB.Create(&newUserProfile).Error; err != nil {
		return "", err
	}

	return newUserProfile.UserId, nil
}

func (s *UserController) IsSocialUserExistAndActive(email string) (string, error) {
	var socialUser dbModel.User

	err := s.RelationalDB.Model(&dbModel.User{}).Where("email = ? AND type = ? OR type = ?", utils.FormatStringToLowerCase(email), constants.APPLE_USER_TYPE, constants.GOOGLE_USER_TYPE).First(&socialUser).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", nil
		}
		return "", err
	}

	return socialUser.UserId, nil
}

// user account deletion operations

func (s *UserController) SoftDeleteUserProfile(userId string, ch chan bool) {
	defer close(ch)

	var userProfile dbModel.User

	if err := s.RelationalDB.Model(&dbModel.User{}).Where("user_id = ?", userId).First(&userProfile).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ch <- true
			return
		}
		fmt.Println(err)
		ch <- false
		return
	}

	res := s.RelationalDB.Model(&dbModel.User{}).Where("user_id = ?", userId).Update("is_delete", true)
	if res.Error != nil || res.RowsAffected == 0 {
		fmt.Println(res.Error)
		ch <- false
		return
	}

	ch <- true
}

func (s *UserController) UpsertDeleteRequet(reason, email, userId string) (bool, error) {
	var existingRequest dbModel.DeleteRequest
	err := s.RelationalDB.Where("user_id = ? AND email = ?", userId, reason).First(&existingRequest).Error

	if err == nil {
		updateData := map[string]interface{}{
			"reason": reason,
		}

		if updateErr := s.RelationalDB.Model(&dbModel.DeleteRequest{}).Where("user_id = ? AND email = ?", userId, email).Updates(updateData).Error; updateErr != nil {
			return false, updateErr
		}
		return true, nil
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		newReq := dbModel.DeleteRequest{
			UserId: userId,
			Reason: reason,
			Email:  email,
		}
		if createErr := s.RelationalDB.Create(&newReq).Error; createErr != nil {
			return false, createErr
		}
		return true, nil
	}

	return false, err
}

func (s *UserController) GenerateAndStoreOtp(userId string, key string) (int, error) {
	var otpStr string

	if os.Getenv("TESTING") == "true" {
		otpStr = constants.TEST_OTP
	} else {
		var err error 
		otpStr, err = utils.GenerateOTP(constants.OTP_LENGTH)
		if err != nil {
			return 0, err
		}
	}

	otpKey := fmt.Sprintf("%s%s", userId, key)
	err := s.RedisDB.Set(ctx, otpKey, otpStr, 15*time.Minute).Err()
	if err != nil {
		return 0, err
	}

	otpInt, err := utils.StringToInt(otpStr)
	if err != nil {
		return 0, err
	}

	return otpInt, nil
}

func (s *UserController) VerifyOTP(otpKey string, otp int) (bool, error) {
	digits := len(strconv.Itoa(otp))
	if digits < constants.OTP_LENGTH {
		return false, errors.New("invalid OTP")
	}

	value, err := s.RedisDB.Get(ctx, otpKey).Result()
	if err == redis.Nil {
		return false, errors.New("OTP expired or type not valid")
	} else if err != nil {
		fmt.Println("error verify otp:", err.Error())
		return false, errors.New("OTP expired or entered invalid OTP")
	}

	if value != fmt.Sprint(otp) {
		return false, errors.New("wrong OTP")
	}

	if _, err := s.RedisDB.Del(ctx, otpKey).Result(); err != nil {
		fmt.Println("error deleting OTP key:", err.Error())
		return false, errors.New("error verifying otp")
	}

	return true, nil
}