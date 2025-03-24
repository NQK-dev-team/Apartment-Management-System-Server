package main

import (
	"api/config"
	"api/routes"
	"api/utils"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	router := gin.New()

	// Load .env file
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	// Middleware
	// CORS
	corsOriginList := config.GetEnv("CORS_ORIGIN")
	if corsOriginList == "" {
		corsOriginList = "*"
	}
	corsOrigin := strings.Split(corsOriginList, ",")
	for v, _ := range corsOrigin {
		corsOrigin[v] = strings.TrimSpace(corsOrigin[v])
	}

	router.Use(cors.New(cors.Config{
		AllowOrigins: corsOrigin,
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		// AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Refresh-Token"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}))
	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	router.Use(gin.Recovery())

	// Serve static files
	router.Static("/public", "./public")

	// Init DB
	err = config.InitDB()
	if err != nil {
		panic(err)
	}

	// Init AWS S3 service
	utils.InitS3Connection()

	// Router settings
	// Limit form size to 50 MB
	router.MaxMultipartMemory = 50 << 20 // 50 MB

	// Init routes
	r := router.Group("api")
	routes.InitRoutes(r)

	// Init custom validation rules
	utils.InitCustomValidationRules()

	// Start server
	var app *http.Server
	isHTTPS := config.GetEnv("HTTPS")
	if isHTTPS == "" {
		isHTTPS = "false"
	}

	port := config.GetEnv("PORT")
	if port == "" {
		port = "8080"
	}

	app = &http.Server{
		Addr:           ":" + port,
		Handler:        router,
		ReadTimeout:    1 * time.Minute,
		WriteTimeout:   1 * time.Minute,
		MaxHeaderBytes: 1 << 20,
	}

	fmt.Println("----------------------------------------")
	fmt.Println("|                                      |")
	fmt.Println("|                                      |")
	fmt.Printf("|   Server is running on port: %s    |\n", port)
	if isHTTPS == "true" {
		fmt.Println("|   TLS is enabled                     |")
	}
	fmt.Println("|                                      |")
	fmt.Println("|                                      |")
	fmt.Println("----------------------------------------")

	certPath := config.GetEnv("CERT_PATH")
	keyPath := config.GetEnv("KEY_PATH")

	if isHTTPS == "true" && certPath != "" && keyPath != "" {
		app.ListenAndServeTLS(certPath, keyPath)
	} else if isHTTPS == "false" {
		app.ListenAndServe()
	}
}
