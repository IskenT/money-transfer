package service

import (
	"context"
	"fmt"
	"time"

	"github.com/IskenT/money-transfer/internal/domain/model"
	"github.com/IskenT/money-transfer/internal/domain/repository"
	"github.com/IskenT/money-transfer/internal/infra/database"
	"github.com/IskenT/money-transfer/internal/infra/repository/postgresql"
	"github.com/jmoiron/sqlx"
)

// TransferService
type TransferService struct {
	userRepo       repository.UserRepository
	transferRepo   repository.TransferRepository
	txManager      *database.TransactionManager
	pgUserRepo     *postgresql.UserRepository
	pgTransferRepo *postgresql.TransferRepository
}

// NewTransferService
func NewTransferService(
	userRepo repository.UserRepository,
	transferRepo repository.TransferRepository,
	txManager *database.TransactionManager,
	pgUserRepo *postgresql.UserRepository,
	pgTransferRepo *postgresql.TransferRepository,
) *TransferService {
	return &TransferService{
		userRepo:       userRepo,
		transferRepo:   transferRepo,
		txManager:      txManager,
		pgUserRepo:     pgUserRepo,
		pgTransferRepo: pgTransferRepo,
	}
}

// CreateTransfer
func (s *TransferService) CreateTransfer(fromUserID, toUserID string, amount int) (*model.Transfer, error) {
	if amount <= 0 {
		return nil, model.ErrInvalidAmount
	}

	if fromUserID == toUserID {
		return nil, model.ErrSameAccount
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var transfer *model.Transfer

	err := s.txManager.WithTransaction(ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		// SELECT FOR UPDATE
		fromUser, err := s.pgUserRepo.GetForUpdate(ctx, tx, fromUserID)
		if err != nil {
			return err
		}

		toUser, err := s.pgUserRepo.GetForUpdate(ctx, tx, toUserID)
		if err != nil {
			return err
		}

		if fromUser.Balance < amount {
			return model.ErrInsufficientFunds
		}

		transferIDGen, err := s.pgTransferRepo.GetTransferIDGenerator()
		if err != nil {
			return err
		}

		txIDGen, err := s.pgTransferRepo.GetTransactionIDGenerator()
		if err != nil {
			return err
		}

		now := time.Now()
		txID := txIDGen()
		stan := model.Stan(txID)

		debitTx := &model.Transaction{
			Stan:            stan,
			Amount:          amount,
			State:           model.TransactionStatePending,
			TransactionType: model.TransactionTypeDebit,
			PaymentSource:   model.PaymentMethodTypeTransfer,
			Note:            fmt.Sprintf("Transfer to %s", toUser.Name),
			CreatedAt:       now,
			UpdatedAt:       now,
		}

		creditTx := &model.Transaction{
			Stan:            stan,
			Amount:          amount,
			State:           model.TransactionStatePending,
			TransactionType: model.TransactionTypeCredit,
			PaymentSource:   model.PaymentMethodTypeTransfer,
			Note:            fmt.Sprintf("Transfer from %s", fromUser.Name),
			CreatedAt:       now,
			UpdatedAt:       now,
		}

		transfer = &model.Transfer{
			ID:         transferIDGen(),
			FromUserID: fromUserID,
			ToUserID:   toUserID,
			Amount:     amount,
			State:      model.TransactionStatePending,
			DebitTx:    debitTx,
			CreditTx:   creditTx,
			CreatedAt:  now,
		}

		fromUser.Balance -= amount
		toUser.Balance += amount

		if err := s.pgUserRepo.UpdateTx(ctx, tx, fromUser); err != nil {
			return err
		}

		if err := s.pgUserRepo.UpdateTx(ctx, tx, toUser); err != nil {
			return err
		}

		transfer.State = model.TransactionStateCompleted
		transfer.DebitTx.State = model.TransactionStateCompleted
		transfer.CreditTx.State = model.TransactionStateCompleted
		transfer.CompletedAt = time.Now()
		transfer.DebitTx.UpdatedAt = transfer.CompletedAt
		transfer.CreditTx.UpdatedAt = transfer.CompletedAt

		return s.pgTransferRepo.CreateTx(ctx, tx, transfer)
	})

	if err != nil {
		return nil, err
	}

	return transfer, nil
}

// GetTransfer
func (s *TransferService) GetTransfer(id string) (*model.Transfer, error) {
	return s.transferRepo.GetByID(id)
}

// ListTransfers
func (s *TransferService) ListTransfers() ([]*model.Transfer, error) {
	return s.transferRepo.List()
}

// ListUsers
func (s *TransferService) ListUsers() ([]*model.User, error) {
	return s.userRepo.List()
}

// UserByID
func (s *TransferService) UserByID(id string) (*model.User, error) {
	return s.userRepo.GetByID(id)
}
