package model

import "time"

// Transfer
type Transfer struct {
	ID          string
	FromUserID  string
	ToUserID    string
	Amount      int
	State       TransactionState
	DebitTx     *Transaction
	CreditTx    *Transaction
	CreatedAt   time.Time
	CompletedAt time.Time
}
