package main

import (
	"api/config"
	"api/constants"
	"api/services"
	"api/utils"
	"fmt"

	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	// Init DB
	err = config.InitDB()
	if err != nil {
		panic(err)
	}

	// Init storage services
	utils.InitStorageServices()

	// Init custom validation rules
	constants.InitCustomValidationRules()

	if config.GetEnv("APM_CLIENT_BASE_URL") == "" {
		fmt.Errorf("APM_CLIENT_BASE_URL is not set")
		return
	}

	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Recovered from panic: %v\n", r)
		}
	}()

	// Run upload cron
	uploadService := services.NewUploadService()
	uploadService.RunUploadCron()

	// if err != nil {
	// 	// return err
	// 	fmt.Printf("Error running master cron: %v\n", err)
	// }
}
