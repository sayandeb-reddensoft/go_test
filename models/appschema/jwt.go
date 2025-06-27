package appschema

// Restructure this data to fit your needs
type JwtData struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Type      int    `json:"type"`
}

// user quick token/temp token claims
type TempJwtData struct {
	Email string `json:"email"`
}
