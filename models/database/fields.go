package model

// user
var USER_COLLECTION_FIELDS = struct {
	ID        string
	USER_ID   string
	EMAIL     string
	PASSWORD  string
	ROLE      string
	IS_ACTIVE string
}{
	ID:        "id",
	USER_ID:   "user_id",
	EMAIL:     "email",
	PASSWORD:  "password",
	ROLE:      "role",
	IS_ACTIVE: "is_active",
}

// role
var ROLE_COLLECTION_FIELDS = struct {
	ID        string
	ROLE_NAME string
}{
	ID:        "id",
	ROLE_NAME: "role_name",
}

// logins
var LOGIN_COLLECTION_FIELDS = struct {
	ID               string
	USER_ID          string
	LOGIN_TIMESTAMP  string
	LOGOUT_TIMESTAMP string
	NO_OF_ATTEMPTS   string
}{
	ID:               "id",
	USER_ID:          "user_id",
	LOGIN_TIMESTAMP:  "login_time_stamp",
	LOGOUT_TIMESTAMP: "log_out_time_stamp",
	NO_OF_ATTEMPTS:   "no_of_attempts",
}

// organization
var ORG_COLLECTION_FIELDS = struct {
	ID       string
	ORG_NAME string
	ORG_USER string
	ADDRESS  string
}{
	ID:       "id",
	ORG_NAME: "org_name",
	ORG_USER: "org_user",
	ADDRESS:  "address",
}

// address
var ADDRESS_COLLECTION_FIELDS = struct {
	ID               string
	STREET_LINE      string
	CITY             string
	STATE            string
	POSTAL_INDEX_CODE string
}{
	ID:               "id",
	STREET_LINE:      "street_line",
	CITY:             "city",
	STATE:            "state",
	POSTAL_INDEX_CODE: "postal_index_code",
}