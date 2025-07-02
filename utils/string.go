package utils

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

func FormatStringToLowerCase(data string) string {
	return strings.ToLower(data)
}

func StringToInt(s string) (int, error) {
	Int, err := strconv.Atoi(s)
	if err != nil {
		fmt.Println("failed to type conversion", s, err)
		return 0, errors.New("can't convert the type string to int")
	}

	return Int, nil
}

func GenerateUserUID() string {
	newUUID := uuid.New()
	key := strings.ReplaceAll(newUUID.String(), "-", "")
	return key
}

