package memory

import (
	"sync"

	"github.com/IskenT/money-transfer/internal/domain/model"
)

// TransferRepository
type TransferRepository struct {
	transfers map[string]*model.Transfer
	mu        sync.RWMutex
}

// NewTransferRepository
func NewTransferRepository() *TransferRepository {
	return &TransferRepository{
		transfers: make(map[string]*model.Transfer),
	}
}

// Create
func (r *TransferRepository) Create(transfer *model.Transfer) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.transfers[transfer.ID] = transfer
	return nil
}

// GetByID
func (r *TransferRepository) GetByID(id string) (*model.Transfer, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	transfer, ok := r.transfers[id]
	if !ok {
		return nil, model.ErrTransferNotFound
	}

	return transfer, nil
}

// List
func (r *TransferRepository) List() ([]*model.Transfer, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	transfers := make([]*model.Transfer, 0, len(r.transfers))
	for _, transfer := range r.transfers {
		transfers = append(transfers, transfer)
	}

	return transfers, nil
}
