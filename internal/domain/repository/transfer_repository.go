package repository

import "github.com/IskenT/money-transfer/internal/domain/model"

// TransferRepository
type TransferRepository interface {
	Create(transfer *model.Transfer) error
	GetByID(id string) (*model.Transfer, error)
	List() ([]*model.Transfer, error)
}
