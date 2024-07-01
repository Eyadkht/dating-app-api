package core

import (
	"dating-app/pkg/models"
	"fmt"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func InitDb() {

	retryAttempts := 3
	retryInterval := 2 * time.Second

	var database *gorm.DB
	var databaseError error

	// Retry connecting to the database in case of failure
	// Docker Compose takes some time to start the database service
	// and the application might start before the database is ready
	// and the connection will fail
	for i := 0; i < retryAttempts; i++ {
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			AppConfig.MYSQL_USER,
			AppConfig.MYSQL_PASSWORD,
			AppConfig.MYSQL_HOST,
			AppConfig.MYSQL_PORT,
			AppConfig.MYSQL_DATABASE,
		)

		database, databaseError = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if databaseError == nil {
			db = database
			break // Connection successful, exit the loop
		}

		log.Printf("Failed to connect to the database (attempt %d/%d): %s", i+1, retryAttempts, databaseError)
		time.Sleep(retryInterval)
	}

	if databaseError != nil {
		log.Fatal("Failed to connect to the database after retries:", databaseError)
	}

	// Migrate models to the Database
	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Token{})
	db.AutoMigrate(&models.Swipe{})
	db.AutoMigrate(&models.Match{})
}

func GetDb() *gorm.DB {
	return db
}
