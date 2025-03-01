package main

import (
	"api/config"
	"api/services"
	"time"

	"github.com/joho/godotenv"
)

func emailWorker() {
	emailService := services.NewEmailQueueService()
	emailService.SendEmail()
}

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	err = config.InitMailer()
	if err != nil {
		panic(err)
	}

	config.InitDB()

	for {
		emailWorker()

		// Sleep for 5 seconds
		time.Sleep(5 * time.Second)
	}
}
