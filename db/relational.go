package db

import (
	"fmt"
	"os"

	"github.com/nelsonin-research-org/cdc-auth/globals"
	model "github.com/nelsonin-research-org/cdc-auth/models/database"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// ConnectPostgres connects to the PostgreSQL database.
func ConnectPostgres(config *gorm.Config) error {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		os.Getenv("RELATIONAL_DB_HOST"),
		os.Getenv("RELATIONAL_DB_USER"),
		os.Getenv("RELATIONAL_DB_PASSWORD"),
		os.Getenv("RELATIONAL_DB_NAME"),
		os.Getenv("RELATIONAL_DB_PORT"),
		os.Getenv("RELATIONAL_SSL_MODE"),
		os.Getenv("RELATIONAL_TIME_ZONE"),
	)

	var err error
	globals.RelationalDb, err = gorm.Open(postgres.Open(dsn), config)
	if err != nil {
		return err
	}

	return relationalSchemaMigrate(globals.RelationalDb)
}

// relationalSchemaMigrate performs schema migration for relational db.
func relationalSchemaMigrate(db *gorm.DB) error {
	if err := db.AutoMigrate(
		&model.Role{},
		&model.User{},
		&model.Login{},
		&model.Address{},
		&model.Organization{},
	); err != nil {
		return err
	}
	fmt.Println("Relational schema migration completed")
	return nil
}
