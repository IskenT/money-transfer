package handler

import (
	"encoding/json"
	"net/http"

	"github.com/IskenT/money-transfer/internal/app/service"
	"github.com/IskenT/money-transfer/internal/domain/model"
	httpModel "github.com/IskenT/money-transfer/internal/infra/http/model"
	"github.com/gorilla/mux"
)

// UserController
type UserController struct {
	service *service.TransferService
}

// NewUserController
func NewUserController(service *service.TransferService) *UserController {
	return &UserController{
		service: service,
	}
}

// GetUserByIDHandler godoc
// @Summary Get a specific user
// @Description Get user details by ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} httpModel.UserResponse
// @Failure 404 {object} httpModel.ErrorResponse
// @Failure 500 {object} httpModel.ErrorResponse
// @Router /api/users/{id} [get]
func (c *UserController) GetUserByIDHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	id := vars["id"]

	user, err := c.service.UserByID(id)
	if err != nil {
		statusCode := http.StatusInternalServerError

		if err == model.ErrUserNotFound {
			statusCode = http.StatusNotFound
		}

		w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(httpModel.ErrorResponse{Error: err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(httpModel.UserToResponse(user))
}

// ListUsersHandler godoc
// @Summary List all users
// @Description Get a list of all users
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {array} httpModel.UserResponse
// @Failure 500 {object} httpModel.ErrorResponse
// @Router /api/users [get]
func (c *UserController) ListUsersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	users, err := c.service.ListUsers()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(httpModel.ErrorResponse{Error: err.Error()})
		return
	}

	response := make([]*httpModel.UserResponse, 0, len(users))
	for _, u := range users {
		response = append(response, httpModel.UserToResponse(u))
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
