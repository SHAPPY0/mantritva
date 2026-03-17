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

	r, cj := api.SetupRouter(db, mongoDB)
	//Start Cron Job
	if cj != nil {
		go cj.Start()
	}

	// wait for interrupt
	// quit := make(chan os.Signal, 1)
	// signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	// <-quit
	// log.Println("Shutting down server...")

	// stop cron and wait for running jobs to finish (with same timeout)
	// if cj != nil {
	// 	stopCtx := cj.Stop()
	// 	<-stopCtx.Done()
	// 	log.Printf("Cron stopped")
	// }

	if err := r.Run(":" + config.GetConfig().GetString("PORT")); err != nil {
		panic(err)
	}
}
