package postgresql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/IskenT/money-transfer/internal/domain/model"
	"github.com/jmoiron/sqlx"
)

// DBUser
type DBUser struct {
	ID        int64     `db:"id"`
	Name      string    `db:"name"`
	Balance   int       `db:"balance"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

// UserRepository
type UserRepository struct {
	db *sqlx.DB
}

// NewUserRepository
func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

// GetByID
func (r *UserRepository) GetByID(id string) (*model.User, error) {
	var dbUser DBUser

	err := r.db.Get(&dbUser, `
		SELECT id, name, balance, created_at, updated_at 
		FROM money_transfer.users 
		WHERE id = $1
	`, id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model.ErrUserNotFound
		}
		return nil, fmt.Errorf("error getting user by ID: %w", err)
	}

	return &model.User{
		ID:      fmt.Sprintf("%d", dbUser.ID),
		Name:    dbUser.Name,
		Balance: dbUser.Balance,
	}, nil
}

// Update
func (r *UserRepository) Update(user *model.User) error {
	_, err := r.db.Exec(`
		UPDATE money_transfer.users 
		SET name = $1, balance = $2, updated_at = NOW() 
		WHERE id = $3
	`, user.Name, user.Balance, user.ID)

	if err != nil {
		return fmt.Errorf("error updating user: %w", err)
	}

	return nil
}

// List
func (r *UserRepository) List() ([]*model.User, error) {
	var dbUsers []DBUser

	err := r.db.Select(&dbUsers, `
		SELECT id, name, balance, created_at, updated_at 
		FROM money_transfer.users 
		ORDER BY id
	`)

	if err != nil {
		return nil, fmt.Errorf("error listing users: %w", err)
	}

	users := make([]*model.User, len(dbUsers))
	for i, dbUser := range dbUsers {
		users[i] = &model.User{
			ID:      fmt.Sprintf("%d", dbUser.ID),
			Name:    dbUser.Name,
			Balance: dbUser.Balance,
		}
	}

	return users, nil
}

// GetForUpdate
func (r *UserRepository) GetForUpdate(ctx context.Context, tx *sqlx.Tx, id string) (*model.User, error) {
	var dbUser DBUser

	err := tx.GetContext(ctx, &dbUser, `
		SELECT id, name, balance, created_at, updated_at 
		FROM money_transfer.users 
		WHERE id = $1
		FOR UPDATE
	`, id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model.ErrUserNotFound
		}
		return nil, fmt.Errorf("error getting user by ID with lock: %w", err)
	}

	return &model.User{
		ID:      fmt.Sprintf("%d", dbUser.ID),
		Name:    dbUser.Name,
		Balance: dbUser.Balance,
	}, nil
}

// UpdateTx
func (r *UserRepository) UpdateTx(ctx context.Context, tx *sqlx.Tx, user *model.User) error {
	_, err := tx.ExecContext(ctx, `
		UPDATE money_transfer.users 
		SET name = $1, balance = $2, updated_at = NOW() 
		WHERE id = $3
	`, user.Name, user.Balance, user.ID)

	if err != nil {
		return fmt.Errorf("error updating user in transaction: %w", err)
	}

	return nil
}
