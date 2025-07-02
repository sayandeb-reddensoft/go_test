package utils

import (
	"os"
)

func IsDevelopment() bool {
	return os.Getenv("ENV_DEV") == "true"
}

func IsStage() bool {
	return os.Getenv("ENV_STAGE") == "true"
}

func IsProduction() bool {
	return os.Getenv("ENV_PROD") == "true"
}