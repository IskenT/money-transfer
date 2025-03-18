package router

import (
	"net/http"

	_ "github.com/IskenT/money-transfer/docs"
	"github.com/IskenT/money-transfer/internal/app/service"
	"github.com/IskenT/money-transfer/internal/infra/http/handler"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

// Router
type Router struct {
	router   *mux.Router
	services *service.Services
}

// NewRouter
func NewRouter(services *service.Services) *Router {
	return &Router{
		router:   mux.NewRouter(),
		services: services,
	}
}

// setupRoutes
func (r *Router) setupRoutes() {
	transferController := handler.NewTransferController(r.services.TransferService)
	userController := handler.NewUserController(r.services.TransferService)

	apiRouter := r.router.PathPrefix("/api").Subrouter()

	apiRouter.HandleFunc("/transfers", transferController.CreateTransferHandler).Methods("POST")
	apiRouter.HandleFunc("/transfers", transferController.ListTransfersHandler).Methods("GET")
	apiRouter.HandleFunc("/transfers/{id}", transferController.GetTransferByIDHandler).Methods("GET")

	apiRouter.HandleFunc("/users", userController.ListUsersHandler).Methods("GET")
	apiRouter.HandleFunc("/users/{id}", userController.GetUserByIDHandler).Methods("GET")

	r.router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	r.router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/swagger/index.html", http.StatusSeeOther)
	})
}

// Handler
func (r *Router) Handler() http.Handler {
	r.setupRoutes()
	return r.router
}
