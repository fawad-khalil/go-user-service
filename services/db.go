package db

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	// DB is the static GORM database connection instance
	DB *gorm.DB
)

func NewSQLiteDB(dbPath string) error {
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return err
	}
	DB = db
	return nil
}

func Close() error {
	db, err := DB.DB()
	if err != nil {
		return err
	}
	return db.Close()
}
