package main

import (
	"github.com/shappy0/mantritva/internal/api"
	"github.com/shappy0/mantritva/pkg/config"
)

func main() {

	//Load environment variables
	config.Load()
	//DB Connection
	db := config.SetupDatabase()
	if db == nil {
		panic("Failed to connect to postgresDB")
	}

	// cassandra, err := config.SetupCassandra()
	// if err != nil {
	// 	panic("Failed to connect to Cassandra")
	// }
	// defer cassandra.Close()
	// cassandra.Ping()

	mongoDB, err := config.SetupMongoDB()
	if err != nil {
		panic(err)
	}

	r := api.SetupRouter(db, mongoDB)

	// Enable CORS middleware
	// r.Use(corsMiddleware())

	if err := r.Run(":8080"); err != nil {
		panic(err)
	}
}
