package repository

import (
	"github.com/IskenT/money-transfer/internal/domain/repository"
	"github.com/IskenT/money-transfer/internal/infra/database"
	"github.com/IskenT/money-transfer/internal/infra/repository/postgresql"
)

// Factory
type Factory struct {
	txManager *database.TransactionManager
}

// NewFactory
func NewFactory(txManager *database.TransactionManager) *Factory {
	return &Factory{
		txManager: txManager,
	}
}

// CreateUserRepository
func (f *Factory) CreateUserRepository() (repository.UserRepository, *postgresql.UserRepository) {
	pgRepo := postgresql.NewUserRepository(f.txManager.DB())
	return pgRepo, pgRepo
}

// CreateTransferRepository
func (f *Factory) CreateTransferRepository() (repository.TransferRepository, *postgresql.TransferRepository) {
	pgRepo := postgresql.NewTransferRepository(f.txManager.DB())
	return pgRepo, pgRepo
}
