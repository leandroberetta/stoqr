package database

import (
	"log"

	"gorm.io/gorm"
)

// Connect connects to a Postgres database if config is ok, otherwise creates a SQLite database
func Connect() *gorm.DB {
	db := &gorm.DB{}
	if err := checkPostgresConfig(); err != nil {
		log.Printf("Falling back to SQLite: %s", err)
		db, err = openSQLite()
		if err != nil {
			log.Fatal(err)
		}
	} else {
		db, err = openPostgres()
		if err != nil {
			log.Fatal(err)
		}
	}
	return db
}
