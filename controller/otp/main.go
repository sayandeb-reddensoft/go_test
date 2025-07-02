package controller

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/nelsonin-research-org/cdc-auth/interfaces"
	limitations "github.com/nelsonin-research-org/cdc-auth/models/limitation"
	"github.com/nelsonin-research-org/cdc-auth/utils"
	"github.com/redis/go-redis/v9"
)


type otpControllerImpl struct {
	Redis   *redis.Client
}

func NewOTPController(r *redis.Client) interfaces.OTPController {
	return &otpControllerImpl{Redis: r}
}

func (s *otpControllerImpl) GenerateOTP(length int) (int, error) {
	otpStr, err := utils.GenerateOTP(length)
	if err != nil {
		return 0, err
	}

	otpInt, err := utils.StringToInt(otpStr)
	if err != nil {
		return 0, err
	}

	return otpInt, nil
}

func (s *otpControllerImpl) StoreOtp(ctx context.Context, ttl time.Duration, key string, otp int) (bool, error) {
	otpKey := fmt.Sprintf("%s", key)

	if err := s.Redis.Set(ctx, otpKey, otp, ttl).Err(); err != nil {
		return false, errors.New("failed to store OTP:" + err.Error())
	}

	return true, nil
}

func (s *otpControllerImpl) VerifyOTP(ctx context.Context, otpKey string, otp int) (bool, error) {
	if len(strconv.Itoa(otp)) < limitations.OTP_LIMITATION.OTP_LENGTH {
		return false, errors.New("OTP is invalid")
	}

	storedOtp, err := s.Redis.Get(ctx, otpKey).Result()
	if err == redis.Nil {
		return false, errors.New("OTP not found or expired")
	}
	if err != nil {
		fmt.Printf("failed to fetch OTP: %w", err)
		return false, errors.New("Internal Server Error")
	}

	if storedOtp != fmt.Sprint(otp) {
		return false, errors.New("incorrect OTP entered")
	}

	if _, err := s.Redis.Del(ctx, otpKey).Result(); err != nil {
		fmt.Printf("failed to delete OTP after verification: %w", err)
		return false, errors.New("Internal Server Error")
	}

	return true, nil
}