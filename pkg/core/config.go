package core

import (
	"os"
)

type Config struct {
	ENVIRONMENT    string
	MYSQL_USER     string
	MYSQL_PASSWORD string
	MYSQL_DATABASE string
	MYSQL_HOST     string
	MYSQL_PORT     string
}

var AppConfig Config

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func LoadConfig() {

	AppConfig = Config{
		ENVIRONMENT:    getEnv("ENVIRONMENT", "development"),
		MYSQL_USER:     getEnv("MYSQL_USER", "muzz_db_user"),
		MYSQL_PASSWORD: getEnv("MYSQL_PASSWORD", "muzz_db_password"),
		MYSQL_DATABASE: getEnv("MYSQL_DATABASE", "muzz_dating_db"),
		MYSQL_HOST:     getEnv("MYSQL_HOST", "db"),
		MYSQL_PORT:     getEnv("MYSQL_PORT", "3306"),
	}
}
