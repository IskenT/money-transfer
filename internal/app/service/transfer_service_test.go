package service_test

import (
	"testing"

	"github.com/IskenT/money-transfer/internal/app/service"
	"github.com/IskenT/money-transfer/internal/domain/model"
	"github.com/IskenT/money-transfer/internal/infra/repository/memory"
)

func TestCreateTransfer(t *testing.T) {
	userRepo := memory.NewUserRepository()
	transferRepo := memory.NewTransferRepository()
	transferService := service.NewTransferService(userRepo, transferRepo)

	markUser, _ := userRepo.GetByID("1")
	janeUser, _ := userRepo.GetByID("2")
	adamUser, _ := userRepo.GetByID("3")

	markInitialBalance := markUser.Balance
	janeInitialBalance := janeUser.Balance
	adamInitialBalance := adamUser.Balance

	tests := []struct {
		name          string
		fromUserID    string
		toUserID      string
		amount        int
		expectedError error
	}{
		{
			name:          "Valid transfer",
			fromUserID:    "1", // Mark
			toUserID:      "2", // Jane
			amount:        1000,
			expectedError: nil,
		},
		{
			name:          "Insufficient funds",
			fromUserID:    "3", // Adam (0 balance)
			toUserID:      "2", // Jane
			amount:        1000,
			expectedError: model.ErrInsufficientFunds,
		},
		{
			name:          "Invalid amount",
			fromUserID:    "1", // Mark
			toUserID:      "2", // Jane
			amount:        0,
			expectedError: model.ErrInvalidAmount,
		},
		{
			name:          "Negative amount",
			fromUserID:    "1", // Mark
			toUserID:      "2", // Jane
			amount:        -100,
			expectedError: model.ErrInvalidAmount,
		},
		{
			name:          "Same account",
			fromUserID:    "1", // Mark
			toUserID:      "1", // Mark
			amount:        1000,
			expectedError: model.ErrSameAccount,
		},
		{
			name:          "User not found",
			fromUserID:    "1", // Mark
			toUserID:      "999",
			amount:        1000,
			expectedError: model.ErrUserNotFound,
		},
	}

	// Run
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			transfer, err := transferService.CreateTransfer(tc.fromUserID, tc.toUserID, tc.amount)

			if tc.expectedError != nil {
				if err != tc.expectedError {
					t.Errorf("Expected error %v, got %v", tc.expectedError, err)
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if transfer.FromUserID != tc.fromUserID {
				t.Errorf("Expected FromUserID %s, got %s", tc.fromUserID, transfer.FromUserID)
			}
			if transfer.ToUserID != tc.toUserID {
				t.Errorf("Expected ToUserID %s, got %s", tc.toUserID, transfer.ToUserID)
			}
			if transfer.Amount != tc.amount {
				t.Errorf("Expected Amount %d, got %d", tc.amount, transfer.Amount)
			}
			if transfer.State != model.TransactionStateCompleted {
				t.Errorf("Expected State %s, got %s", model.TransactionStateCompleted, transfer.State)
			}

			fromUser, _ := userRepo.GetByID(tc.fromUserID)
			toUser, _ := userRepo.GetByID(tc.toUserID)

			expectedFromBalance := 0
			expectedToBalance := 0

			switch tc.fromUserID {
			case "1":
				expectedFromBalance = markInitialBalance - tc.amount
			case "2":
				expectedFromBalance = janeInitialBalance - tc.amount
			case "3":
				expectedFromBalance = adamInitialBalance - tc.amount
			}

			switch tc.toUserID {
			case "1":
				expectedToBalance = markInitialBalance + tc.amount
			case "2":
				expectedToBalance = janeInitialBalance + tc.amount
			case "3":
				expectedToBalance = adamInitialBalance + tc.amount
			}

			if fromUser.Balance != expectedFromBalance {
				t.Errorf("Expected FromUser balance %d, got %d", expectedFromBalance, fromUser.Balance)
			}
			if toUser.Balance != expectedToBalance {
				t.Errorf("Expected ToUser balance %d, got %d", expectedToBalance, toUser.Balance)
			}
		})
	}
}

func TestConcurrentTransfers(t *testing.T) {
	userRepo := memory.NewUserRepository()
	transferRepo := memory.NewTransferRepository()
	transferService := service.NewTransferService(userRepo, transferRepo)

	markUser, _ := userRepo.GetByID("1")
	janeUser, _ := userRepo.GetByID("2")
	initialMarkBalance := markUser.Balance
	initialJaneBalance := janeUser.Balance

	numTransfers := 10
	transferAmount := 100

	done := make(chan bool, numTransfers*2)

	for i := 0; i < numTransfers; i++ {
		go func() {
			_, err := transferService.CreateTransfer("1", "2", transferAmount)
			if err != nil {
				t.Errorf("Error in Mark to Jane transfer: %v", err)
			}
			done <- true
		}()

		go func() {
			_, err := transferService.CreateTransfer("2", "1", transferAmount)
			if err != nil {
				t.Errorf("Error in Jane to Mark transfer: %v", err)
			}
			done <- true
		}()
	}

	for i := 0; i < numTransfers*2; i++ {
		<-done
	}
	markUser, _ = userRepo.GetByID("1")
	janeUser, _ = userRepo.GetByID("2")

	if markUser.Balance != initialMarkBalance {
		t.Errorf("Expected Mark's final balance to be %d, got %d",
			initialMarkBalance, markUser.Balance)
	}

	if janeUser.Balance != initialJaneBalance {
		t.Errorf("Expected Jane's final balance to be %d, got %d",
			initialJaneBalance, janeUser.Balance)
	}
}
