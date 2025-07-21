package main

import (
	"api/config"
	"api/migrations"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

var userMigration *migrations.UserMigration
var refreshTokenMigration *migrations.RefreshTokenMigration
var emailVerifyTokenMigration *migrations.EmailVerifyTokenMigration
var passwordResetTokenMigration *migrations.PasswordResetTokenMigration
var buildingMigration *migrations.BuildingMigration
var contractMigration *migrations.ContractMigration

// var messageMigration *migrations.MessageMigration
var notificationMigration *migrations.NotificationMigration
var supportTicketMigration *migrations.SupportTicketMigration
var emailQueueMigration *migrations.EmailQueueMigration

func migrateUp() {
	fmt.Println("--------- Migrate Up Start ---------")
	userMigration.Up()
	refreshTokenMigration.Up()
	emailVerifyTokenMigration.Up()
	passwordResetTokenMigration.Up()
	buildingMigration.Up()
	contractMigration.Up()
	// messageMigration.Up()
	notificationMigration.Up()
	supportTicketMigration.Up()
	emailQueueMigration.Up()
	fmt.Println("--------- Migrate Up Finish ---------")
}

func migrateDown() {
	fmt.Println("--------- Migrating Down Start ---------")
	supportTicketMigration.Down()
	notificationMigration.Down()
	// messageMigration.Down()
	contractMigration.Down()
	buildingMigration.Down()
	passwordResetTokenMigration.Down()
	emailVerifyTokenMigration.Down()
	refreshTokenMigration.Down()
	userMigration.Down()
	emailQueueMigration.Down()
	fmt.Println("--------- Migrating Down Finish ---------")
}

func initMigrations() {
	userMigration = migrations.NewUserMigration()
	refreshTokenMigration = migrations.NewRefreshTokenMigration()
	emailVerifyTokenMigration = migrations.NewEmailVerifyTokenMigration()
	passwordResetTokenMigration = migrations.NewPasswordResetTokenMigration()
	buildingMigration = migrations.NewBuildingMigration()
	contractMigration = migrations.NewContractMigration()
	// messageMigration = migrations.NewMessageMigration()
	notificationMigration = migrations.NewNotificationMigration()
	supportTicketMigration = migrations.NewSupportTicketMigration()
	emailQueueMigration = migrations.NewEmailQueueMigration()
}

func migrateHandler(mode string) {
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
	initMigrations()
	switch mode {
	case "up":
		migrateUp()
	case "down":
		migrateDown()
	default:
		fmt.Println("Invalid mode. Use 'up' or 'down'.")
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please provide a migration mode: up or down")
		return
	}
	fmt.Println("Migration mode:", os.Args[1])
	migrateHandler(os.Args[1])
}
