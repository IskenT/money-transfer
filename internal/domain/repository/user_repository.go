package repository

import "github.com/IskenT/money-transfer/internal/domain/model"

// UserRepository
type UserRepository interface {
	GetByID(id string) (*model.User, error)
	Update(user *model.User) error
	List() ([]*model.User, error)
}
