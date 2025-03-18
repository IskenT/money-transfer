package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/IskenT/money-transfer/internal/config"
	"github.com/IskenT/money-transfer/internal/infra/database"
	migrate "github.com/rubenv/sql-migrate"
)

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("Usage: migrate [up|down]")
		os.Exit(1)
	}

	direction := args[0]
	if direction != "up" && direction != "down" {
		fmt.Println("Direction must be either 'up' or 'down'")
		os.Exit(1)
	}

	cfg := config.NewConfig()
	dbConfig := database.NewDBConfig(cfg)

	db, err := database.NewDBWithRetry(dbConfig, 5, 3*time.Second)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	migrations := &migrate.FileMigrationSource{
		Dir: "migrations",
	}

	var n int
	if strings.EqualFold(direction, "up") {
		n, err = migrate.Exec(db.DB, "postgres", migrations, migrate.Up)
	} else {
		n, err = migrate.Exec(db.DB, "postgres", migrations, migrate.Down)
	}

	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	fmt.Printf("Applied %d migrations %s\n", n, direction)
}
