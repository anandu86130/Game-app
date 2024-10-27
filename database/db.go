package db

import (
	"gaming/model"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DB is a global variable for the database connection
var DB *gorm.DB

// DBconnect initializes the database connection
func DBconnect() {
	// Fetching the Data Source Name (DSN) from the environment variables
	dsn := os.Getenv("DSN")

	// Connecting to the database
	var err error
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	//error handling
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	DB = db
	// Auto-migrate the database schema
	err = DB.AutoMigrate(&model.User{}, &model.OTP{}, &model.League{}, &model.Team{}, &model.Tournament{}, &model.TeamA{}, &model.TeamB{}) // Specify models to migrate
	//error handling
	if err != nil {
		log.Fatal("failed to auto migrate", err)
	}
}
