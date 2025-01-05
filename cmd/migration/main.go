package main

import (
	"api/config"
	"api/migrations"
	"fmt"
	"os"
)

var userMigration *migrations.UserMigration
var refreshTokenMigration *migrations.RefreshTokenMigration
var emailVerifyTokenMigration *migrations.EmailVerifyTokenMigration
var passwordResetTokenMigration *migrations.PasswordResetTokenMigration
var buildingMigration *migrations.BuildingMigration
var roomMigration *migrations.RoomMigration

func migrateUp() {
	fmt.Println("--------- Migrate Up Start ---------")
	userMigration.Up()
	refreshTokenMigration.Up()
	emailVerifyTokenMigration.Up()
	passwordResetTokenMigration.Up()
	buildingMigration.Up()
	roomMigration.Up()
	fmt.Println("--------- Migrate Up Finish ---------")
}

func migrateDown() {
	fmt.Println("--------- Migrating Down Start ---------")
	userMigration.Down()
	refreshTokenMigration.Down()
	emailVerifyTokenMigration.Down()
	passwordResetTokenMigration.Down()
	buildingMigration.Down()
	roomMigration.Down()
	fmt.Println("--------- Migrating Down Finish ---------")
}

func initMigrations() {
	userMigration = migrations.NewUserMigration()
	refreshTokenMigration = migrations.NewRefreshTokenMigration()
	emailVerifyTokenMigration = migrations.NewEmailVerifyTokenMigration()
	passwordResetTokenMigration = migrations.NewPasswordResetTokenMigration()
	buildingMigration = migrations.NewBuildingMigration()
	roomMigration = migrations.NewRoomMigration()
}

func migrateHandler(mode string) {
	// Init DB
	err := config.InitDB()
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
