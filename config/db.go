package config

import (
	"fmt"

	postgresDriver "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB() error {
	host, _ := GetEnv("DB_HOST")
	port, _ := GetEnv("DB_PORT")
	user, _ := GetEnv("DB_USER")
	password, _ := GetEnv("DB_PASSWORD")
	dbName, _ := GetEnv("DB_NAME")
	ssh, _ := GetEnv("DB_SSH")
	timeZone, _ := GetEnv("APP_TIME_ZONE")
	appEnv, _ := GetEnv("APP_ENV")

	var logMode logger.LogLevel

	if appEnv == "production" {
		logMode = logger.Silent
	} else {
		logMode = logger.Info
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s timeZone=%s", host, port, user, password, dbName, ssh, timeZone)

	var err error

	DB, err = gorm.Open(postgresDriver.Open(dsn), &gorm.Config{
		Logger:      logger.Default.LogMode(logMode),
		PrepareStmt: true,
	})

	if err != nil {
		return err
	}

	return nil
}
