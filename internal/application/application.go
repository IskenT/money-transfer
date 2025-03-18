package application

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/IskenT/money-transfer/internal/app/service"
	"github.com/IskenT/money-transfer/internal/config"
	"github.com/IskenT/money-transfer/internal/infra/database"
	"github.com/IskenT/money-transfer/internal/infra/http/middleware"
	"github.com/IskenT/money-transfer/internal/infra/http/router"
	repository "github.com/IskenT/money-transfer/internal/infra/repository/factory"
	"github.com/jmoiron/sqlx"
)

// Application represents the main application
type Application struct {
	server    *http.Server
	services  *service.Services
	router    *router.Router
	isRunning bool
	db        *sqlx.DB
	txManager *database.TransactionManager
}

// NewApplication
func NewApplication() *Application {
	cfg := config.NewConfig()

	dbConfig := database.NewDBConfig(cfg)
	db, err := database.NewDBWithRetry(dbConfig, 5, 3*time.Second)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	txManager := database.NewTransactionManager(db)

	repoFactory := repository.NewFactory(txManager)

	userRepo, pgUserRepo := repoFactory.CreateUserRepository()
	transferRepo, pgTransferRepo := repoFactory.CreateTransferRepository()

	transferService := service.NewTransferService(
		userRepo, transferRepo, txManager, pgUserRepo, pgTransferRepo,
	)

	services := &service.Services{
		TransferService: transferService,
	}

	r := router.NewRouter(services)

	server := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      middleware.ApplyCORS(r.Handler()),
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	return &Application{
		server:    server,
		services:  services,
		router:    r,
		isRunning: false,
		db:        db,
		txManager: txManager,
	}
}

// Start
func (a *Application) Start() error {
	if a.isRunning {
		return fmt.Errorf("server is already running")
	}

	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, syscall.SIGINT, syscall.SIGTERM)
		sig := <-sigint
		log.Printf("Received signal: %v. Shutting down...", sig)

		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		if err := a.server.Shutdown(ctx); err != nil {
			log.Printf("HTTP server shutdown error: %v", err)
		}

		if err := a.db.Close(); err != nil {
			log.Printf("Database connection close error: %v", err)
		}

		log.Println("Server gracefully stopped")
		close(idleConnsClosed)
	}()

	a.isRunning = true
	log.Printf("Server started on http://localhost%s", a.server.Addr)
	log.Printf("Swagger UI available at http://localhost%s/swagger/index.html", a.server.Addr)

	if err := a.server.ListenAndServe(); err != http.ErrServerClosed {
		return fmt.Errorf("HTTP server error: %v", err)
	}

	<-idleConnsClosed
	return nil
}

// Stop
func (a *Application) Stop() error {
	if !a.isRunning {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	a.isRunning = false

	if err := a.server.Shutdown(ctx); err != nil {
		return err
	}

	if err := a.db.Close(); err != nil {
		return err
	}

	return nil
}

// DB
func (a *Application) DB() *sqlx.DB {
	return a.db
}
