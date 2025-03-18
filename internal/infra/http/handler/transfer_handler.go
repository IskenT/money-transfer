package handler

import (
	"encoding/json"
	"net/http"

	"github.com/IskenT/money-transfer/internal/app/service"
	"github.com/IskenT/money-transfer/internal/domain/model"
	httpModel "github.com/IskenT/money-transfer/internal/infra/http/model"
	"github.com/gorilla/mux"
)

// TransferController handles HTTP requests for transfers
type TransferController struct {
	service *service.TransferService
}

// NewTransferController creates a new TransferController
func NewTransferController(service *service.TransferService) *TransferController {
	return &TransferController{
		service: service,
	}
}

// CreateTransferHandler godoc
// @Summary Create a new money transfer
// @Description Transfer money from one user to another
// @Tags transfers
// @Accept json
// @Produce json
// @Param transfer body httpModel.TransferRequest true "Transfer details"
// @Success 201 {object} httpModel.TransferResponse
// @Failure 400 {object} httpModel.ErrorResponse
// @Failure 404 {object} httpModel.ErrorResponse
// @Failure 500 {object} httpModel.ErrorResponse
// @Router /api/transfers [post]
func (c *TransferController) CreateTransferHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req httpModel.TransferRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(httpModel.ErrorResponse{Error: "Invalid request format"})
		return
	}

	transfer, err := c.service.CreateTransfer(req.FromUserID, req.ToUserID, req.Amount)
	if err != nil {
		statusCode := http.StatusInternalServerError

		switch err {
		case model.ErrInsufficientFunds:
			statusCode = http.StatusBadRequest
		case model.ErrSameAccount:
			statusCode = http.StatusBadRequest
		case model.ErrUserNotFound:
			statusCode = http.StatusNotFound
		case model.ErrInvalidAmount:
			statusCode = http.StatusBadRequest
		}

		w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(httpModel.ErrorResponse{Error: err.Error()})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(httpModel.TransferToResponse(transfer))
}

// GetTransferByIDHandler godoc
// @Summary Get a specific transfer
// @Description Get transfer details by ID
// @Tags transfers
// @Accept json
// @Produce json
// @Param id path string true "Transfer ID"
// @Success 200 {object} httpModel.TransferResponse
// @Failure 404 {object} httpModel.ErrorResponse
// @Failure 500 {object} httpModel.ErrorResponse
// @Router /api/transfers/{id} [get]
func (c *TransferController) GetTransferByIDHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	id := vars["id"]

	transfer, err := c.service.GetTransfer(id)
	if err != nil {
		statusCode := http.StatusInternalServerError

		if err == model.ErrTransferNotFound {
			statusCode = http.StatusNotFound
		}

		w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(httpModel.ErrorResponse{Error: err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(httpModel.TransferToResponse(transfer))
}

// ListTransfersHandler godoc
// @Summary List all transfers
// @Description Get a list of all transfers
// @Tags transfers
// @Accept json
// @Produce json
// @Success 200 {array} httpModel.TransferResponse
// @Failure 500 {object} httpModel.ErrorResponse
// @Router /api/transfers [get]
func (c *TransferController) ListTransfersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	transfers, err := c.service.ListTransfers()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(httpModel.ErrorResponse{Error: err.Error()})
		return
	}

	// Convert to response objects
	response := make([]*httpModel.TransferResponse, 0, len(transfers))
	for _, t := range transfers {
		response = append(response, httpModel.TransferToResponse(t))
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
