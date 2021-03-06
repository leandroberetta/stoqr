package database

import (
	"errors"
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func checkPostgresConfig() error {
	if host := os.Getenv("STOQR_API_DB_HOST"); host == "" {
		return errors.New(("postgres host is empty"))
	}
	if port := os.Getenv("STOQR_API_DB_PORT"); port == "" {
		return errors.New(("postgres port is empty"))
	}
	if user := os.Getenv("STOQR_API_DB_USER"); user == "" {
		return errors.New(("postgres user is empty"))
	}
	if password := os.Getenv("STOQR_API_DB_PASSWORD"); password == "" {
		return errors.New(("postgres password is empty"))
	}
	if name := os.Getenv("STOQR_API_DB_NAME"); name == "" {
		return errors.New(("postgres db name is empty"))
	}
	return nil
}

func generatePostgresConnectionString() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("STOQR_API_DB_HOST"),
		os.Getenv("STOQR_API_DB_USER"),
		os.Getenv("STOQR_API_DB_PASSWORD"),
		os.Getenv("STOQR_API_DB_NAME"),
		os.Getenv("STOQR_API_DB_PORT"))
}

func openPostgres() (*gorm.DB, error) {
	return gorm.Open(postgres.Open(generatePostgresConnectionString()), &gorm.Config{})
}
