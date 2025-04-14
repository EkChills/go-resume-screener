package database

import (
	"fmt"

	"github.com/ekchills/go-resume-screener/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDb() error {
	dsn := "host=localhost user=damned password=12345678 dbname=postgres port=5432 sslmode=disable"
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	fmt.Println("Connected to database")
	DB = db
	return nil
}

func MigrateDb() error {
	// Migrate the schema
	if err := DB.AutoMigrate(&models.User{}, &models.Resume{}); err != nil {
		return err
	}
	return nil
}