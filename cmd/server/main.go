package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/IskenT/money-transfer/internal/app/service"
	"github.com/IskenT/money-transfer/internal/infra/http/handler"
	"github.com/IskenT/money-transfer/internal/infra/repository/memory"

	_ "github.com/IskenT/money-transfer/docs"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {
	// repositories
	userRepo := memory.NewUserRepository()
	transferRepo := memory.NewTransferRepository()

	// service
	transferService := service.NewTransferService(userRepo, transferRepo)

	// controller
	transferController := handler.NewTransferController(transferService)
	userController := handler.NewUserController(transferService)

	// router
	router := mux.NewRouter()

	apiRouter := router.PathPrefix("/api").Subrouter()

	// endpoints
	apiRouter.HandleFunc("/transfers", transferController.CreateTransferHandler).Methods("POST")
	apiRouter.HandleFunc("/transfers", transferController.ListTransfersHandler).Methods("GET")
	apiRouter.HandleFunc("/transfers/{id}", transferController.GetTransferByIDHandler).Methods("GET")

	// endpoints
	apiRouter.HandleFunc("/users", userController.ListUsersHandler).Methods("GET")
	apiRouter.HandleFunc("/users/{id}", userController.GetUserByIDHandler).Methods("GET")

	// Swagger UI
	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/swagger/index.html", http.StatusSeeOther)
	})

	// Apply CORS middleware
	corsMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}

	handler := corsMiddleware(router)

	port := ":8080"
	fmt.Printf("Server started on http://localhost%s\n", port)
	fmt.Printf("Swagger UI available at http://localhost%s/swagger/index.html\n", port)
	log.Fatal(http.ListenAndServe(port, handler))
}
