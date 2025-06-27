package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

// GenerateOTP generates a OTP on given length
func GenerateOTP(length int) (string, error) {
	max := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(length)), nil) 

	randNum, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}

	otp := fmt.Sprintf("%0*s", length, randNum.String()) 

	return otp, nil
}