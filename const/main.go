package constants

var DB_DRIVERS = []string{"postgres", "redis"}
var EMAIL_USER_TYPE = 1

// from the request body
var RESET_REQUEST_TYPE = "reset"
var VERIFY_REQUEST_TYPE = "verify"

// in redis
var FORGET_KEY = "_forget:"
var VERIFY_KEY = "_signup:"

// otp
var OTP_LENGTH = 4

// token types
var LONG_LIVE_TOKEN = 1
var TEMP_TOKEN = 2
var REFRESH_TOKEN = 3

// formats and responses
var DATE_FORMAT = "2006-01-02" // date format in YYYY-MM-DD