package main

import (
	"api/config"
	"api/services"
	"api/utils"
	"sync"
	"time"

	"github.com/joho/godotenv"
)

const NUMBER_OF_ROUTINES int = 4 // Number of goroutines to run concurrently
var emailService *services.EmailQueueService
var contractService *services.ContractService
var roomService *services.RoomService
var billService *services.BillService

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

func billWorker(wg *sync.WaitGroup, momoCounter *int) {
	defer wg.Done() // Decrement the counter when this goroutine finishes
	*momoCounter++
	if (*momoCounter)%3 == 0 {
		billService.GetMomoResult()
		*momoCounter = 0
	}
	billService.UpdateBillStatus()
}

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	// Init mailer
	err = config.InitMailer()
	if err != nil {
		panic(err)
	}

	// Init DB
	config.InitDB()

	// Init MoMo config
	utils.InitMoMoConfig()

	emailService = services.NewEmailQueueService()
	contractService = services.NewContractService()
	roomService = services.NewRoomService()
	billService = services.NewBillService()

	// for {
	// 	emailWorker()

	// 	// Sleep for 5 seconds
	// 	time.Sleep(5 * time.Second)
	// }

	// Counter to keep track of how many times the loop has run, when the momoCounter reaches 3, execute getMomoResult inside billWorker
	momoCounter := 0

	// Create a ticker that ticks every 10 seconds
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop() // Ensure the ticker is stopped when main exits

	// Use a WaitGroup to keep track of active goroutines
	var wg sync.WaitGroup

	for range ticker.C {
		wg.Add(NUMBER_OF_ROUTINES)
		go emailWorker(&wg)
		go contractWorker(&wg)
		go roomWorker(&wg)
		go billWorker(&wg, &momoCounter)
		wg.Wait() // Wait for all workers to finish before continuing to the next tick
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
