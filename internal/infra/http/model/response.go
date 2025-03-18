package model

import (
	"fmt"
	"time"

	domainModel "github.com/IskenT/money-transfer/internal/domain/model"
)

// ErrorResponse
type ErrorResponse struct {
	Error string `json:"error" example:"insufficient funds"`
}

// UserResponse
type UserResponse struct {
	ID               string `json:"id" example:"1"`
	Name             string `json:"name" example:"Mark"`
	Balance          int    `json:"balance" example:"10000"`
	BalanceFormatted string `json:"balance_formatted" example:"$100.00"`
}

// TransactionResponse
type TransactionResponse struct {
	Stan            string `json:"stan" example:"TRX1647881234567"`
	Amount          int    `json:"amount" example:"1000"`
	AmountFormatted string `json:"amount_formatted" example:"$10.00"`
	State           string `json:"state" example:"COMPLETED"`
	TransactionType string `json:"transaction_type" example:"DEBIT"`
	PaymentSource   string `json:"payment_source" example:"TRANSFER"`
	Note            string `json:"note" example:"Transfer to Jane"`
	CreatedAt       string `json:"created_at" example:"2023-04-10T12:34:56Z"`
	UpdatedAt       string `json:"updated_at" example:"2023-04-10T12:34:56Z"`
}

// TransferResponse
type TransferResponse struct {
	ID              string               `json:"id" example:"TRF1647881234567"`
	FromUserID      string               `json:"from_user_id" example:"1"`
	ToUserID        string               `json:"to_user_id" example:"2"`
	Amount          int                  `json:"amount" example:"1000"`
	AmountFormatted string               `json:"amount_formatted" example:"$10.00"`
	State           string               `json:"state" example:"COMPLETED"`
	DebitTx         *TransactionResponse `json:"debit_tx,omitempty"`
	CreditTx        *TransactionResponse `json:"credit_tx,omitempty"`
	CreatedAt       string               `json:"created_at" example:"2023-04-10T12:34:56Z"`
	CompletedAt     string               `json:"completed_at,omitempty" example:"2023-04-10T12:34:56Z"`
}

// TransferRequest
type TransferRequest struct {
	FromUserID string `json:"from_user_id" example:"1" description:"ID of the sender"`
	ToUserID   string `json:"to_user_id" example:"2" description:"ID of the recipient"`
	Amount     int    `json:"amount" example:"1000" description:"Amount to transfer in cents (e.g., 1000 = $10.00)"`
}

// FormatMoney
func FormatMoney(cents int) string {
	dollars := float64(cents) / 100.0
	return fmt.Sprintf("$%.2f", dollars)
}

// FormatTime
func FormatTime(t time.Time) string {
	return t.Format(time.RFC3339)
}

// TransferToResponse
func TransferToResponse(t *domainModel.Transfer) *TransferResponse {
	res := &TransferResponse{
		ID:              t.ID,
		FromUserID:      t.FromUserID,
		ToUserID:        t.ToUserID,
		Amount:          t.Amount,
		AmountFormatted: FormatMoney(t.Amount),
		State:           string(t.State),
		CreatedAt:       FormatTime(t.CreatedAt),
	}

	if !t.CompletedAt.IsZero() {
		completed := FormatTime(t.CompletedAt)
		res.CompletedAt = completed
	}

	if t.DebitTx != nil {
		res.DebitTx = &TransactionResponse{
			Stan:            string(t.DebitTx.Stan),
			Amount:          t.DebitTx.Amount,
			AmountFormatted: FormatMoney(t.DebitTx.Amount),
			State:           string(t.DebitTx.State),
			TransactionType: string(t.DebitTx.TransactionType),
			PaymentSource:   string(t.DebitTx.PaymentSource),
			Note:            t.DebitTx.Note,
			CreatedAt:       FormatTime(t.DebitTx.CreatedAt),
			UpdatedAt:       FormatTime(t.DebitTx.UpdatedAt),
		}
	}

	if t.CreditTx != nil {
		res.CreditTx = &TransactionResponse{
			Stan:            string(t.CreditTx.Stan),
			Amount:          t.CreditTx.Amount,
			AmountFormatted: FormatMoney(t.CreditTx.Amount),
			State:           string(t.CreditTx.State),
			TransactionType: string(t.CreditTx.TransactionType),
			PaymentSource:   string(t.CreditTx.PaymentSource),
			Note:            t.CreditTx.Note,
			CreatedAt:       FormatTime(t.CreditTx.CreatedAt),
			UpdatedAt:       FormatTime(t.CreditTx.UpdatedAt),
		}
	}

	return res
}

// UserToResponse
func UserToResponse(u *domainModel.User) *UserResponse {
	return &UserResponse{
		ID:               u.ID,
		Name:             u.Name,
		Balance:          u.Balance,
		BalanceFormatted: FormatMoney(u.Balance),
	}
}
