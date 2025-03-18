package model

import "time"

// Stan
type Stan string

// TransactionType
type TransactionType string

// TransactionState
type TransactionState string

// PaymentMethodType d
type PaymentMethodType string

const (
	TransactionTypeDebit  TransactionType = "DEBIT"
	TransactionTypeCredit TransactionType = "CREDIT"

	TransactionStatePending   TransactionState = "PENDING"
	TransactionStateCompleted TransactionState = "COMPLETED"
	TransactionStateFailed    TransactionState = "FAILED"

	PaymentMethodTypeTransfer PaymentMethodType = "TRANSFER"
)

// Transaction
type Transaction struct {
	Stan            Stan
	AppID           int
	ProfileID       uint32
	CardID          uint32
	InstanceID      string
	Card            string
	CardName        string
	CardAcct        string
	BankC           string
	Expiry          string
	Amount          int
	Fee             int
	State           TransactionState
	TransactionType TransactionType
	PaymentSource   PaymentMethodType
	Merchant        string
	Note            string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
