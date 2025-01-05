package main

import (
	"api/config"
	"api/services"
	"fmt"
)

func main() {
	if appEnv, _ := config.GetEnv("APP_ENV"); appEnv == "production" {
		fmt.Println("This command is only available in development environment.")
		return
	}

	var table int

	fmt.Println("--------- Import into table ---------")
	fmt.Println("1 - User")
	fmt.Println("2 - Building")
	fmt.Println("3 - Contract")
	fmt.Println("4 - Bill")
	fmt.Print("Your choice: ")
	fmt.Scanln(&table)

	// Init DB
	err := config.InitDB()
	if err != nil {
		panic(err)
	}
	var importSerive = services.NewImportServiceDev()

	fmt.Println("--------- Importing File Start ---------")
	if err := importSerive.ImportFileDev(table); err != nil {
		fmt.Println("Error importing file: ", err)
		return
	}
	fmt.Println("--------- Importing File Finish ---------")
}
