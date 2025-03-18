package service

import (
	"fmt"
	"sync"
	"time"

	"github.com/IskenT/money-transfer/internal/domain/model"
	"github.com/IskenT/money-transfer/internal/domain/repository"
)

// TransferService
type TransferService struct {
	userRepo     repository.UserRepository
	transferRepo repository.TransferRepository
	mu           sync.Mutex
}

// NewTransferService
func NewTransferService(userRepo repository.UserRepository, transferRepo repository.TransferRepository) *TransferService {
	return &TransferService{
		userRepo:     userRepo,
		transferRepo: transferRepo,
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

	s.mu.Lock()
	defer s.mu.Unlock()

	// Get users
	fromUser, err := s.userRepo.GetByID(fromUserID)
	if err != nil {
		return nil, err
	}

	toUser, err := s.userRepo.GetByID(toUserID)
	if err != nil {
		return nil, err
	}

	// Check balance
	if fromUser.Balance < amount {
		return nil, model.ErrInsufficientFunds
	}

	now := time.Now()
	stan := model.Stan(fmt.Sprintf("TRX%d", time.Now().UnixNano()))

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

	transfer := &model.Transfer{
		ID:         fmt.Sprintf("TRF%d", time.Now().UnixNano()),
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

	if err := s.userRepo.Update(fromUser); err != nil {
		return nil, err
	}

	if err := s.userRepo.Update(toUser); err != nil {
		fromUser.Balance += amount
		s.userRepo.Update(fromUser)
		return nil, err
	}

	transfer.State = model.TransactionStateCompleted
	transfer.DebitTx.State = model.TransactionStateCompleted
	transfer.CreditTx.State = model.TransactionStateCompleted
	transfer.CompletedAt = time.Now()
	transfer.DebitTx.UpdatedAt = transfer.CompletedAt
	transfer.CreditTx.UpdatedAt = transfer.CompletedAt

	if err := s.transferRepo.Create(transfer); err != nil {
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
