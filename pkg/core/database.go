package core

import (
	"fmt"
	"log"
	"muzz-dating/pkg/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func InitDb() {

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		AppConfig.MYSQL_USER,
		AppConfig.MYSQL_PASSWORD,
		AppConfig.MYSQL_HOST,
		AppConfig.MYSQL_PORT,
		AppConfig.MYSQL_DATABASE,
	)
	database, database_error := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	db = database

	if database_error != nil {
		log.Fatal("Failed to connect to the database:", database_error)
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
