package model

type LoginData struct {
	Email    string `json:"email"     validate:"required,email,max=255"`
	Password string `json:"password" validate:"required"`
}

type SignUpData struct {
	FirstName string `json:"first_name" validate:"required,max=255"`
	LastName  string `json:"last_name"  validate:"max=255"`
	DOB       string `json:"dob"    validate:"required"`
	Email     string `json:"email"     validate:"required,email,max=255"`
	Password  string `json:"password"  validate:"required,min=8,max=1000"`
}

type ForgetPasswordData struct {
	Email string `json:"email" validate:"required"`
}

type VerifyOtpData struct {
	Email string `json:"email" validate:"required"`
	Type  string `json:"type" validate:"required"`
	OTP   int    `json:"otp" validate:"required"`
}

type ResendOtpData struct {
	Email string `json:"email" validate:"required"`
	Type  string `json:"type" validate:"required"`
}

type GetUserData struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Gender    string `json:"gender"`
	Email     string `json:"email"`
	Type      int    `json:"type"`
	DOB       string `json:"dob"`
	IsActive  bool   `json:"is_active"`
}

type UpdatePasswordData struct {
	Password string `json:"password"  validate:"required,min=8,max=1000"`
}
