package database

import (
	"log"

	"github.com/leandroberetta/stoqr/stoqr-api/models"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Connect() *gorm.DB {
	db := &gorm.DB{}
	if err := checkPostgresConfig(); err != nil {
		log.Printf("Falling back to SQLite: %s", err)
		db, err = gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
		if err != nil {
			log.Fatal(err)
		}
	} else {
		db, err = gorm.Open(postgres.Open(generatePostgresConnectionString()), &gorm.Config{})
		if err != nil {
			log.Fatal(err)
		}
	}
	db.AutoMigrate(&models.Item{})
	return db
}
