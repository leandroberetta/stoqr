package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func openSQLite() (*gorm.DB, error) {
	return gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
}
