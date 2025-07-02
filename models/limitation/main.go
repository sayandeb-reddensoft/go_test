package model

import "time"

// request limitors
var HANDLER_LIMITATION = struct {
	UPDATE_PASSWORD  int
	FORGOT_PASSWORD  int
	RESEND_OTP       int
}{
	UPDATE_PASSWORD: 5,
	FORGOT_PASSWORD: 10,
	RESEND_OTP:      5,
}

// otp limitation
var OTP_LIMITATION = struct {
	OTP_LENGTH  int
	OTP_EXP     time.Duration 	
}{
	OTP_LENGTH: 4,
	OTP_EXP:    15*time.Minute,
}

// JWT limitation
var JWT_LIMITATION = struct {
	PRIMARY_TOKEN_TTL time.Duration
	REFRESH_TOKEN_TTL time.Duration
	TEMP_TOKEN_TTL    time.Duration
}{
	PRIMARY_TOKEN_TTL: 14*24*time.Hour,
	REFRESH_TOKEN_TTL: 21*24*time.Hour,
	TEMP_TOKEN_TTL:    5*time.Minute,
}