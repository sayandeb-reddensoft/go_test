package utils

import (
	"fmt"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

// Password must be 8 characters long, must include minimum one uppercase, one lowercase, one number and one special character
func IsStrongPassword(password string) bool {
	if len(password) < 8 {
		return false
	}

	hasUppercase := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLowercase := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasDigit := regexp.MustCompile(`\d`).MatchString(password)
	hasSpecialChar := regexp.MustCompile(`[!"#$%&'()*+,-./:;<=>?@[\]^_{|}~]`).MatchString(password)

	return hasUppercase && hasLowercase && hasDigit && hasSpecialChar
}

// HashPassword is used to encrypt the password
func HashPassword(password string) string {
	if password == "" {
		return ""
	}

	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		fmt.Println("failed to hash the password", err)
		return "failed to hash the password"
	}

	return string(bytes)
}

// VerifyPassword is used to verify the password
func VerifyPassword(userPassword string, providedPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(userPassword), []byte(providedPassword))
	check := true
	msg := ""

	if err != nil {
		msg = "wrong password"
		check = false
	}

	return check, msg
}

// Return true if password matches with hashedPassword
func MatchWithHashPassword(password string, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
