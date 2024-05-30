package core

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func InitDb() {

	MYSQL_USER := "muzz_db_user"
	MYSQL_PASSWORD := "muzz_db_password"
	MYSQL_DATABASE := "muzz_dating_db"
	MYSQL_HOST := "db"
	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		MYSQL_USER,
		MYSQL_PASSWORD,
		MYSQL_HOST,
		MYSQL_DATABASE)
	database, database_error := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	db = database

	if database_error != nil {
		log.Fatal("Failed to connect to the database:", database_error)
	}
}

func GetDb() *gorm.DB {
	return db
}
