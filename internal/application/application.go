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
	"github.com/IskenT/money-transfer/internal/infra/http/middleware"
	"github.com/IskenT/money-transfer/internal/infra/http/router"
	"github.com/IskenT/money-transfer/internal/infra/repository/memory"
)

// Application
type Application struct {
	server    *http.Server
	services  *service.Services
	router    *router.Router
	isRunning bool
}

// NewApplication
func NewApplication() *Application {
	userRepo := memory.NewUserRepository()
	transferRepo := memory.NewTransferRepository()

	transferService := service.NewTransferService(userRepo, transferRepo)
	services := &service.Services{
		TransferService: transferService,
	}

	r := router.NewRouter(services)

	server := &http.Server{
		Addr:    ":8080",
		Handler: middleware.ApplyCORS(r.Handler()),
	}

	return &Application{
		server:    server,
		services:  services,
		router:    r,
		isRunning: false,
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
	return a.server.Shutdown(ctx)
}
