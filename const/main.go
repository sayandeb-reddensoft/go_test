package constants

// user & roles
type UserRole int
const (
	SuperAdmin UserRole = 1
	OrgAdmin   UserRole = 2
	Doctor 	   UserRole = 3
	Therapist  UserRole = 4
	Parent     UserRole = 5
)

// database
var DB_DRIVERS = []string{"postgres", "redis"}

// token types
var PRIMARY_TOKEN = 1
var REFRESH_TOKEN = 2
var TEMP_TOKEN    = 3

// formats and responses
var DATE_FORMAT = "2006-01-02" // date format in YYYY-MM-DD

// email & otp
var (
	EMAIL_DELETE_ACCOUNT_SUBJECT = "Account Deletion Requested"
	EMAIL_ONBOARD_ACCOUNT_SUBJECT = "Welcome to Clenz!"
	EMAIL_FORGOT_PASSWORD_ACCOUNT_SUBJECT = "Reset Password"
)

var (
	EMAIL_DELETE_ACCOUNT_TEMPLATE = "delete-account.html"
	EMAIL_ONBOARD_ACCOUNT_TEMPLATE = "welcome.html"
	EMAIL_FORGOT_PASSWORD_ACCOUNT_TEMPLATE = "forgotpassword.html"
)

var (
	RESET_REQUEST_TYPE = "reset"
	VERIFY_REQUEST_TYPE = "verify"

	MEMORY_FORGET_KEY = "_forget:"
	MEMORY_VERIFY_KEY = "_signup:"
) 
