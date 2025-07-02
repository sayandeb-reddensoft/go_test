package model

type SignUpData struct {
	OrgName            string        `json:"org_name" validate:"required"`
	OrgEmail           string        `json:"email" validate:"required,email,max=255"`
	Password           string        `json:"password" validate:"required"`
	ContactName        string        `json:"contact_name" validate:"required"`
	Address            SignupAddress `json:"address" validate:"required"`
	ContactNumber      string        `json:"contact_number" validate:"required"`
	ContactDescription string        `json:"contact_description" validate:"required"`
}

type SignupAddress struct {
	StreetLane string `json:"street_line" validate:"required"`
	City       string `json:"city"  validate:"required"`
	State      string `json:"state" validate:"required"`
	PostalCode string `json:"zip_code" validate:"required"`
}

type User struct {
	UserId   string `json:"user_id"`
	Password string `json:"passwprd"`
	Email    string `json:"email"`
	Role     int    `json:"role_id"`
	IsActive bool   `json:"is_active"`
}

type LoginData struct {
	Email    string `json:"email"     validate:"required,email,max=255"`
	Password string `json:"password" validate:"required"`
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
	Name     string `json:"name"`
	Email    string `json:"email"`
	Role     int    `json:"role"`
	IsActive bool   `json:"is_active"`
}

type UpdatePasswordData struct {
	Password string `json:"password"  validate:"required,min=8,max=1000"`
}
