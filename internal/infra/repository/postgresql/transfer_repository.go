package postgresql

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/IskenT/money-transfer/internal/domain/model"
	"github.com/jmoiron/sqlx"
)

// DBTransaction
type DBTransaction struct {
	ID              int64     `db:"id"`
	Stan            string    `db:"stan"`
	Amount          int       `db:"amount"`
	State           string    `db:"state"`
	TransactionType string    `db:"transaction_type"`
	PaymentSource   string    `db:"payment_source"`
	Note            string    `db:"note"`
	CreatedAt       time.Time `db:"created_at"`
	UpdatedAt       time.Time `db:"updated_at"`
}

// DBTransfer
type DBTransfer struct {
	ID           int64         `db:"id"`
	TransferCode string        `db:"transfer_code"`
	FromUserID   int64         `db:"from_user_id"`
	ToUserID     int64         `db:"to_user_id"`
	Amount       int           `db:"amount"`
	State        string        `db:"state"`
	DebitTxID    sql.NullInt64 `db:"debit_tx_id"`
	CreditTxID   sql.NullInt64 `db:"credit_tx_id"`
	CreatedAt    time.Time     `db:"created_at"`
	CompletedAt  sql.NullTime  `db:"completed_at"`
}

// DBOutboxEvent
type DBOutboxEvent struct {
	ID            int64        `db:"id"`
	AggregateType string       `db:"aggregate_type"`
	AggregateID   string       `db:"aggregate_id"`
	EventType     string       `db:"event_type"`
	Payload       []byte       `db:"payload"`
	CreatedAt     time.Time    `db:"created_at"`
	ProcessedAt   sql.NullTime `db:"processed_at"`
}

// TransferRepository
type TransferRepository struct {
	db *sqlx.DB
}

// NewTransferRepository
func NewTransferRepository(db *sqlx.DB) *TransferRepository {
	return &TransferRepository{
		db: db,
	}
}

// Create
func (r *TransferRepository) Create(transfer *model.Transfer) error {
	ctx := context.Background()
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("error beginning transaction: %w", err)
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	err = r.CreateTx(ctx, tx, transfer)
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("error committing transaction: %w", err)
	}

	return nil
}

// CreateTx
func (r *TransferRepository) CreateTx(ctx context.Context, tx *sqlx.Tx, transfer *model.Transfer) error {
	var debitTxID int64
	err := tx.QueryRowxContext(ctx, `
		INSERT INTO money_transfer.transactions (
			stan, amount, state, transaction_type, payment_source, note, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8
		) RETURNING id
	`,
		transfer.DebitTx.Stan,
		transfer.DebitTx.Amount,
		transfer.DebitTx.State,
		transfer.DebitTx.TransactionType,
		transfer.DebitTx.PaymentSource,
		transfer.DebitTx.Note,
		transfer.DebitTx.CreatedAt,
		transfer.DebitTx.UpdatedAt,
	).Scan(&debitTxID)

	if err != nil {
		return fmt.Errorf("error inserting debit transaction: %w", err)
	}

	var creditTxID int64
	err = tx.QueryRowxContext(ctx, `
		INSERT INTO money_transfer.transactions (
			stan, amount, state, transaction_type, payment_source, note, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8
		) RETURNING id
	`,
		transfer.CreditTx.Stan,
		transfer.CreditTx.Amount,
		transfer.CreditTx.State,
		transfer.CreditTx.TransactionType,
		transfer.CreditTx.PaymentSource,
		transfer.CreditTx.Note,
		transfer.CreditTx.CreatedAt,
		transfer.CreditTx.UpdatedAt,
	).Scan(&creditTxID)

	if err != nil {
		return fmt.Errorf("error inserting credit transaction: %w", err)
	}

	completedAt := sql.NullTime{}
	if !transfer.CompletedAt.IsZero() {
		completedAt.Time = transfer.CompletedAt
		completedAt.Valid = true
	}

	var transferID int64
	err = tx.QueryRowxContext(ctx, `
		INSERT INTO money_transfer.transfers (
			transfer_code, from_user_id, to_user_id, amount, state, 
			debit_tx_id, credit_tx_id, created_at, completed_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9
		) RETURNING id
	`,
		transfer.ID,
		transfer.FromUserID,
		transfer.ToUserID,
		transfer.Amount,
		transfer.State,
		debitTxID,
		creditTxID,
		transfer.CreatedAt,
		completedAt,
	).Scan(&transferID)

	if err != nil {
		return fmt.Errorf("error inserting transfer: %w", err)
	}

	payload, err := json.Marshal(map[string]interface{}{
		"transfer_id":  transfer.ID,
		"from_user_id": transfer.FromUserID,
		"to_user_id":   transfer.ToUserID,
		"amount":       transfer.Amount,
		"state":        transfer.State,
		"created_at":   transfer.CreatedAt,
		"completed_at": transfer.CompletedAt,
	})

	if err != nil {
		return fmt.Errorf("error marshaling outbox event payload: %w", err)
	}

	_, err = tx.ExecContext(ctx, `
		INSERT INTO money_transfer.outbox_events (
			aggregate_type, aggregate_id, event_type, payload
		) VALUES (
			'transfer', $1, 'transfer_completed', $2
		)
	`, transfer.ID, payload)

	if err != nil {
		return fmt.Errorf("error inserting outbox event: %w", err)
	}

	return nil
}

// GetByID
func (r *TransferRepository) GetByID(id string) (*model.Transfer, error) {
	var dbTransfer DBTransfer

	err := r.db.Get(&dbTransfer, `
		SELECT id, transfer_code, from_user_id, to_user_id, amount, state, 
		       debit_tx_id, credit_tx_id, created_at, completed_at
		FROM money_transfer.transfers
		WHERE transfer_code = $1
	`, id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model.ErrTransferNotFound
		}
		return nil, fmt.Errorf("error getting transfer by ID: %w", err)
	}

	var debitTx DBTransaction
	if dbTransfer.DebitTxID.Valid {
		err = r.db.Get(&debitTx, `
			SELECT id, stan, amount, state, transaction_type, payment_source, note, created_at, updated_at
			FROM money_transfer.transactions
			WHERE id = $1
		`, dbTransfer.DebitTxID.Int64)

		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("error getting debit transaction: %w", err)
		}
	}

	var creditTx DBTransaction
	if dbTransfer.CreditTxID.Valid {
		err = r.db.Get(&creditTx, `
			SELECT id, stan, amount, state, transaction_type, payment_source, note, created_at, updated_at
			FROM money_transfer.transactions
			WHERE id = $1
		`, dbTransfer.CreditTxID.Int64)

		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("error getting credit transaction: %w", err)
		}
	}

	transfer := &model.Transfer{
		ID:         dbTransfer.TransferCode,
		FromUserID: fmt.Sprintf("%d", dbTransfer.FromUserID),
		ToUserID:   fmt.Sprintf("%d", dbTransfer.ToUserID),
		Amount:     dbTransfer.Amount,
		State:      model.TransactionState(dbTransfer.State),
		CreatedAt:  dbTransfer.CreatedAt,
	}

	if dbTransfer.CompletedAt.Valid {
		transfer.CompletedAt = dbTransfer.CompletedAt.Time
	}

	if dbTransfer.DebitTxID.Valid {
		transfer.DebitTx = &model.Transaction{
			Stan:            model.Stan(debitTx.Stan),
			Amount:          debitTx.Amount,
			State:           model.TransactionState(debitTx.State),
			TransactionType: model.TransactionType(debitTx.TransactionType),
			PaymentSource:   model.PaymentMethodType(debitTx.PaymentSource),
			Note:            debitTx.Note,
			CreatedAt:       debitTx.CreatedAt,
			UpdatedAt:       debitTx.UpdatedAt,
		}
	}

	if dbTransfer.CreditTxID.Valid {
		transfer.CreditTx = &model.Transaction{
			Stan:            model.Stan(creditTx.Stan),
			Amount:          creditTx.Amount,
			State:           model.TransactionState(creditTx.State),
			TransactionType: model.TransactionType(creditTx.TransactionType),
			PaymentSource:   model.PaymentMethodType(creditTx.PaymentSource),
			Note:            creditTx.Note,
			CreatedAt:       creditTx.CreatedAt,
			UpdatedAt:       creditTx.UpdatedAt,
		}
	}

	return transfer, nil
}

// List
func (r *TransferRepository) List() ([]*model.Transfer, error) {
	var dbTransfers []DBTransfer

	err := r.db.Select(&dbTransfers, `
		SELECT id, transfer_code, from_user_id, to_user_id, amount, state, 
		       debit_tx_id, credit_tx_id, created_at, completed_at
		FROM money_transfer.transfers
		ORDER BY created_at DESC
	`)

	if err != nil {
		return nil, fmt.Errorf("error listing transfers: %w", err)
	}

	var txIDs []int64
	for _, t := range dbTransfers {
		if t.DebitTxID.Valid {
			txIDs = append(txIDs, t.DebitTxID.Int64)
		}
		if t.CreditTxID.Valid {
			txIDs = append(txIDs, t.CreditTxID.Int64)
		}
	}

	if len(txIDs) == 0 {
		transfers := make([]*model.Transfer, len(dbTransfers))
		for i, dbT := range dbTransfers {
			transfers[i] = &model.Transfer{
				ID:         dbT.TransferCode,
				FromUserID: fmt.Sprintf("%d", dbT.FromUserID),
				ToUserID:   fmt.Sprintf("%d", dbT.ToUserID),
				Amount:     dbT.Amount,
				State:      model.TransactionState(dbT.State),
				CreatedAt:  dbT.CreatedAt,
			}

			if dbT.CompletedAt.Valid {
				transfers[i].CompletedAt = dbT.CompletedAt.Time
			}
		}
		return transfers, nil
	}

	query, args, err := sqlx.In(`
		SELECT id, stan, amount, state, transaction_type, payment_source, note, created_at, updated_at
		FROM money_transfer.transactions
		WHERE id IN (?)
	`, txIDs)

	if err != nil {
		return nil, fmt.Errorf("error preparing transaction query: %w", err)
	}

	query = r.db.Rebind(query)
	var dbTransactions []DBTransaction
	err = r.db.Select(&dbTransactions, query, args...)

	if err != nil {
		return nil, fmt.Errorf("error getting transactions: %w", err)
	}

	txMap := make(map[int64]*model.Transaction)
	for _, tx := range dbTransactions {
		txMap[tx.ID] = &model.Transaction{
			Stan:            model.Stan(tx.Stan),
			Amount:          tx.Amount,
			State:           model.TransactionState(tx.State),
			TransactionType: model.TransactionType(tx.TransactionType),
			PaymentSource:   model.PaymentMethodType(tx.PaymentSource),
			Note:            tx.Note,
			CreatedAt:       tx.CreatedAt,
			UpdatedAt:       tx.UpdatedAt,
		}
	}

	transfers := make([]*model.Transfer, len(dbTransfers))
	for i, dbT := range dbTransfers {
		transfers[i] = &model.Transfer{
			ID:         dbT.TransferCode,
			FromUserID: fmt.Sprintf("%d", dbT.FromUserID),
			ToUserID:   fmt.Sprintf("%d", dbT.ToUserID),
			Amount:     dbT.Amount,
			State:      model.TransactionState(dbT.State),
			CreatedAt:  dbT.CreatedAt,
		}

		if dbT.CompletedAt.Valid {
			transfers[i].CompletedAt = dbT.CompletedAt.Time
		}

		if dbT.DebitTxID.Valid {
			if tx, ok := txMap[dbT.DebitTxID.Int64]; ok {
				transfers[i].DebitTx = tx
			}
		}

		if dbT.CreditTxID.Valid {
			if tx, ok := txMap[dbT.CreditTxID.Int64]; ok {
				transfers[i].CreditTx = tx
			}
		}
	}

	return transfers, nil
}

// GetTransferIDGenerator
func (r *TransferRepository) GetTransferIDGenerator() (func() string, error) {
	return func() string {
		var nextID int64
		_ = r.db.Get(&nextID, `
			SELECT nextval('money_transfer.transfers_id_seq')
		`)
		return fmt.Sprintf("TRF%d", nextID)
	}, nil
}

// GetTransactionIDGenerator
func (r *TransferRepository) GetTransactionIDGenerator() (func() string, error) {
	return func() string {
		var nextID int64
		_ = r.db.Get(&nextID, `
			SELECT nextval('money_transfer.transactions_id_seq')
		`)
		return fmt.Sprintf("TRX%d", nextID)
	}, nil
}
