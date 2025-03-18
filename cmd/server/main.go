package main

import (
	"fmt"
	"log"
	"os"

	_ "github.com/IskenT/money-transfer/docs"
	"github.com/IskenT/money-transfer/internal/app/service/processor"
	"github.com/IskenT/money-transfer/internal/application"
)

// @title Money Transfer API
// @version 1.0
// @description API for a concurrent money transfer system with PostgreSQL storage
// @BasePath /api
func main() {
	app := application.NewApplication()

	outboxProcessor := processor.NewOutboxProcessor(app.DB())
	outboxProcessor.Start()

	defer func() {
		outboxProcessor.Stop()

		if err := app.Stop(); err != nil {
			fmt.Fprintf(os.Stderr, "Error stopping application: %v\n", err)
		}
	}()

	if err := app.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "Error starting server: %v\n", err)
		log.Fatal(err)
	}
}
