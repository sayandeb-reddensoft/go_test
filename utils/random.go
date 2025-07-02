package utils

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
)

// GenerateOTP generates a OTP on given length
func GenerateOTP(length int) (string, error) {
	if length <= 0 {
		return "", errors.New("invalid OTP length")
	}

	firstDigitMax := big.NewInt(9)
	firstDigit, err := rand.Int(rand.Reader, firstDigitMax)
	if err != nil {
		return "", err
	}
	firstDigitStr := fmt.Sprintf("%d", firstDigit.Int64()+1) 

	remainingLength := length - 1
	if remainingLength == 0 {
		return firstDigitStr, nil
	}

	max := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(remainingLength)), nil)
	randNum, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}
	remainingDigits := fmt.Sprintf("%0*s", remainingLength, randNum.String())

	return firstDigitStr + remainingDigits, nil
}