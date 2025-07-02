package main

import (
	"flag"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	constants "github.com/nelsonin-research-org/cdc-auth/const"
	"github.com/nelsonin-research-org/cdc-auth/db"
	"github.com/nelsonin-research-org/cdc-auth/globals"
	app "github.com/nelsonin-research-org/cdc-auth/handlers/data"

	"github.com/nelsonin-research-org/cdc-auth/middleware"
	"github.com/nelsonin-research-org/cdc-auth/routes"
	"github.com/nelsonin-research-org/cdc-auth/utils"
)

var router *gin.Engine

func Setup() error {
	flag.Parse()
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Load .env
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file: %v", err)
		return err
	}

	// Load keys
	if utils.IsDevelopment() {
		err := utils.LoadCertificateKeys()
		if err != nil {
			log.Printf("Error loading certificate keys: %v", err)
			return err
		}
	} 

	// Connect to all DB drivers
	drivers := constants.DB_DRIVERS
	for _, driver := range drivers {
		err := db.ConnectDB(driver)
		if err != nil { 
			return err
		}
	}

	return nil
}

func main() {
	if err := Setup(); err != nil {
		log.Fatalf("Error initializing application: %s", err)
	}

	defer db.DisconnectDB()
 
	globals.RequestStore.Requests = make(map[string]map[string]int)
	handlers := app.LoadAppHandlers()

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(utils.GetCorsConfig())

	versionRouter := router.Group("/api/v1")
	{
		// No auth routes
		routes.NoAuthGroupRoutes(versionRouter, handlers)
	}

	router.NoRoute(middleware.PathNotFound())

	if utils.IsDevelopment(){
		log.Println("Starting Server Locally...")
		router.Run(":8080")
	} else {
		log.Println("Env not specified")
	}
}
