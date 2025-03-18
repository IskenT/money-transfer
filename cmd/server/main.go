package main

import (
	"fmt"
	"log"
	"os"

	_ "github.com/IskenT/money-transfer/docs"
	"github.com/IskenT/money-transfer/internal/application"
)

// @title Money Transfer API
// @version 1.0
// @description API for a concurrent money transfer system
// @BasePath /api
func main() {
	app := application.NewApplication()

	if err := app.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "Error starting server: %v\n", err)
		log.Fatal(err)
	}
}
