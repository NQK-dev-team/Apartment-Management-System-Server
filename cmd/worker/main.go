package main

import (
	"api/config"
	"api/services"
	"time"
)

func emailWorker() {
	emailService := services.NewEmailQueueService()
	emailService.SendEmail()
}

func main() {
	config.InitMailer()
	config.InitDB()

	for {
		emailWorker()

		// Sleep for 5 seconds
		time.Sleep(5 * time.Second)
	}
}
