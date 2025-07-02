package services

import (
	"context"
	"time"

	"github.com/nelsonin-research-org/cdc-auth/interfaces"
)

type otpServiceImpl struct {
	otpController interfaces.OTPController
}

func NewOTPService(c interfaces.OTPController) interfaces.OTPService {
	return &otpServiceImpl{otpController: c}
}

func (service *otpServiceImpl) GenerateOTP(length int) (int, error) {
	return service.otpController.GenerateOTP(length)
}

func (service *otpServiceImpl) StoreOtp(ctx context.Context, ttl time.Duration, otpKey string, otp int) (bool, error) {
	return service.otpController.StoreOtp(ctx, ttl, otpKey, otp)
}

func (service *otpServiceImpl) VerifyOTP(ctx context.Context, otpKey string, otp int) (bool, error) {
	return service.otpController.VerifyOTP(ctx, otpKey, otp)
}

