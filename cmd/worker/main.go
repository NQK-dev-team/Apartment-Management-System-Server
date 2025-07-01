package main

import (
	"api/config"
	"api/services"
	"sync"
	"time"

	"github.com/joho/godotenv"
)

const NUMBER_OF_ROUTINES int = 3 // Number of goroutines to run concurrently
var emailService *services.EmailQueueService
var contractService *services.ContractService
var roomService *services.RoomService

func emailWorker(wg *sync.WaitGroup) {
	defer wg.Done() // Decrement the counter when this goroutine finishes
	emailService.SendEmail()
}

func contractWorker(wg *sync.WaitGroup) {
	defer wg.Done() // Decrement the counter when this goroutine finishes
	contractService.UpdateContractStatus()
}

func roomWorker(wg *sync.WaitGroup) {
	defer wg.Done() // Decrement the counter when this goroutine finishes
	roomService.UpdateRoomStatus()
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

	emailService = services.NewEmailQueueService()
	contractService = services.NewContractService()
	roomService = services.NewRoomService()

	// for {
	// 	emailWorker()

	// 	// Sleep for 5 seconds
	// 	time.Sleep(5 * time.Second)
	// }

	// Create a ticker that ticks every 5 seconds
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop() // Ensure the ticker is stopped when main exits

	// Use a WaitGroup to keep track of active goroutines
	var wg sync.WaitGroup

	for range ticker.C {
		wg.Add(NUMBER_OF_ROUTINES)
		go emailWorker(&wg)
		go contractWorker(&wg)
		go roomWorker(&wg)
		// wg.Wait() // Wait for all workers to finish before continuing to the next tick
	}

	// Important: In a real long-running service, you would typically listen for OS signals (like Ctrl+C)
	// to trigger a graceful shutdown. When a shutdown signal is received, you would:
	// 1. Stop accepting new tasks (e.g., break out of the ticker loop).
	// 2. Call wg.Wait() to ensure all currently running goroutines complete their work.
	// 3. Then exit the main function.

	// As your current `for range ticker.C` loop runs indefinitely, `wg.Wait()` here
	// would never be reached under normal operation. If you ever break out of that loop,
	// then `wg.Wait()` would be useful.
	// For testing and simple observation, you might temporarily make the loop run a fixed number of times.

	// Example of how you *would* wait if the loop was finite:
	// log.Println("Waiting for all active workers to finish...")
	// wg.Wait() // This would block main until all goroutines added to wg have called Done()
	// log.Println("All workers finished. Application shutting down.")
}
