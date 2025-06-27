package utils

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	constants "github.com/nelsonin-research-org/clenz-auth/const"
)

func FormatStringToLowerCase(data string) string {
	return strings.ToLower(data)
}

func FormatStringToDate(date string) (*time.Time, error) {
	if date == "" {
		return nil, nil
	}

	dob, err := time.Parse(constants.DATE_FORMAT, date)
	if err != nil {
		return nil, err
	}

	return &dob, nil
}

func FormatDateToString(dob *time.Time) string {
	if dob == nil {
		return ""
	}
	return dob.Format(constants.DATE_FORMAT)
}

// convert the string to int
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

