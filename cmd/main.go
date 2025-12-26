package main

import (
	"log"

	"github.com/shappy0/mantritva/internal/api"
	"github.com/shappy0/mantritva/pkg/config"
)

func main() {

	//Load environment variables
	config.Load()
	log.Println("Starting Service...")
	log.Println(config.GetConfig().GetString("DB_HOST"))
	//DB Connection
	db := config.SetupDatabase()
	if db == nil {
		panic("Failed to connect to postgresDB")
	}

	mongoDB, err := config.SetupMongoDB()
	if err != nil {
		panic(err.Error())
	}

	r := api.SetupRouter(db, mongoDB)

	// Enable CORS middleware
	// r.Use(corsMiddleware())

	if err := r.Run(":8080"); err != nil {
		panic(err)
	}
}
