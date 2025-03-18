package memory_test

import (
	"testing"
	"time"

	"github.com/IskenT/money-transfer/internal/domain/model"
	"github.com/IskenT/money-transfer/internal/infra/repository/memory"
)

func TestTransferRepository(t *testing.T) {
	transferRepo := memory.NewTransferRepository()

	now := time.Now()
	sampleTransfer := &model.Transfer{
		ID:         "TRF123456789",
		FromUserID: "1",
		ToUserID:   "2",
		Amount:     1000,
		State:      model.TransactionStateCompleted,
		DebitTx: &model.Transaction{
			Stan:            "TRX123456789",
			Amount:          1000,
			State:           model.TransactionStateCompleted,
			TransactionType: model.TransactionTypeDebit,
			PaymentSource:   model.PaymentMethodTypeTransfer,
			Note:            "Test debit",
			CreatedAt:       now,
			UpdatedAt:       now,
		},
		CreditTx: &model.Transaction{
			Stan:            "TRX123456789",
			Amount:          1000,
			State:           model.TransactionStateCompleted,
			TransactionType: model.TransactionTypeCredit,
			PaymentSource:   model.PaymentMethodTypeTransfer,
			Note:            "Test credit",
			CreatedAt:       now,
			UpdatedAt:       now,
		},
		CreatedAt:   now,
		CompletedAt: now,
	}

	t.Run("Create transfer", func(t *testing.T) {
		err := transferRepo.Create(sampleTransfer)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
	})

	t.Run("GetByID existing transfer", func(t *testing.T) {
		transfer, err := transferRepo.GetByID("TRF123456789")
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if transfer.ID != "TRF123456789" {
			t.Errorf("Expected transfer ID TRF123456789, got %s", transfer.ID)
		}

		if transfer.Amount != 1000 {
			t.Errorf("Expected amount 1000, got %d", transfer.Amount)
		}

		if transfer.FromUserID != "1" || transfer.ToUserID != "2" {
			t.Errorf("Expected transfer from user 1 to user 2, got from %s to %s",
				transfer.FromUserID, transfer.ToUserID)
		}
	})

	t.Run("GetByID non-existing transfer", func(t *testing.T) {
		_, err := transferRepo.GetByID("non-existent-id")
		if err != model.ErrTransferNotFound {
			t.Errorf("Expected error %v, got %v", model.ErrTransferNotFound, err)
		}
	})

	t.Run("List transfers", func(t *testing.T) {
		anotherTransfer := &model.Transfer{
			ID:         "TRF987654321",
			FromUserID: "2",
			ToUserID:   "1",
			Amount:     500,
			State:      model.TransactionStateCompleted,
			DebitTx: &model.Transaction{
				Stan:            "TRX987654321",
				Amount:          500,
				State:           model.TransactionStateCompleted,
				TransactionType: model.TransactionTypeDebit,
				PaymentSource:   model.PaymentMethodTypeTransfer,
				Note:            "Another test debit",
				CreatedAt:       now,
				UpdatedAt:       now,
			},
			CreditTx: &model.Transaction{
				Stan:            "TRX987654321",
				Amount:          500,
				State:           model.TransactionStateCompleted,
				TransactionType: model.TransactionTypeCredit,
				PaymentSource:   model.PaymentMethodTypeTransfer,
				Note:            "Another test credit",
				CreatedAt:       now,
				UpdatedAt:       now,
			},
			CreatedAt:   now,
			CompletedAt: now,
		}

		transferRepo.Create(anotherTransfer)

		transfers, err := transferRepo.List()
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if len(transfers) != 2 {
			t.Errorf("Expected 2 transfers, got %d", len(transfers))
		}

		transferMap := make(map[string]*model.Transfer)
		for _, tr := range transfers {
			transferMap[tr.ID] = tr
		}

		if _, exists := transferMap["TRF123456789"]; !exists {
			t.Errorf("Expected transfer with ID TRF123456789 to exist")
		}

		if _, exists := transferMap["TRF987654321"]; !exists {
			t.Errorf("Expected transfer with ID TRF987654321 to exist")
		}
	})
}
