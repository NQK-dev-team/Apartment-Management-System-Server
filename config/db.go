package config

import (
	"fmt"

	postgresDriver "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB() error {
	host := GetEnv("DB_HOST")
	port := GetEnv("DB_PORT")
	user := GetEnv("DB_USER")
	password := GetEnv("DB_PASSWORD")
	dbName := GetEnv("DB_NAME")
	ssh := GetEnv("DB_SSH")
	timeZone := GetEnv("APP_TIME_ZONE")
	// appEnv := GetEnv("APP_ENV")

	// var logMode logger.LogLevel

	// if appEnv == "production" {
	// 	logMode = logger.Silent
	// } else {
	// 	logMode = logger.Info
	// }

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s timeZone=%s", host, port, user, password, dbName, ssh, timeZone)

	var err error

	DB, err = gorm.Open(postgresDriver.Open(dsn), &gorm.Config{
		Logger:      NewModelFileLogger("assets/logs", logger.Info),
		PrepareStmt: true,
	})

	if err != nil {
		return err
	}

	return nil
}
