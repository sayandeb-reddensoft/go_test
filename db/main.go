package db

import (
	"fmt"

	"github.com/nelsonin-research-org/clenz-auth/globals"
	"gorm.io/gorm"
)

// ConnectDB connects to the app database based on the driver.
func ConnectDB(dbDriver string) error {
	var err error

	switch dbDriver {
	case "postgres":
		config := &gorm.Config{}
		err = ConnectPostgres(config)
	case "redis":
		globals.RedisClient, err = ConnectRedis()
	default:
		err = fmt.Errorf("unsupported database driver: %s", dbDriver)
	}

	return err
}

// DisconnectDB close all the connections.
func DisconnectDB() {
	if globals.RelationalDb != nil {
		dbSQL, err := globals.RelationalDb.DB()
		if err != nil {
			panic(err)
		}
		dbSQL.Close()
		fmt.Println("Postgres disconnected")
	}

	if globals.RedisClient != nil {
		err := globals.RedisClient.Close()
		if err != nil {
			fmt.Println("Error disconnecting Redis:", err)
		} else {
			fmt.Println("Redis disconnected")
		}
	}
}
