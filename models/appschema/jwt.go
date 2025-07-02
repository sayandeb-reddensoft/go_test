package appschema

// Restructure this data to fit your needs
type JwtData struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Role      int    `json:"role"`
}

// user quick token/temp token claims
type TempJwtData struct {
	Email string `json:"email"`
}
