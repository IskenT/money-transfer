package memory

import (
	"sync"

	"github.com/IskenT/money-transfer/internal/domain/model"
)

// UserRepository
type UserRepository struct {
	users map[string]*model.User
	mu    sync.RWMutex
}

// NewUserRepository
func NewUserRepository() *UserRepository {
	users := map[string]*model.User{
		"1": {ID: "1", Name: "Mark", Balance: 10000},
		"2": {ID: "2", Name: "Jane", Balance: 5000},
		"3": {ID: "3", Name: "Adam", Balance: 0},
	}

	return &UserRepository{
		users: users,
	}
}

// GetByID
func (r *UserRepository) GetByID(id string) (*model.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, ok := r.users[id]
	if !ok {
		return nil, model.ErrUserNotFound
	}

	return user, nil
}

// Update
func (r *UserRepository) Update(user *model.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.users[user.ID]; !ok {
		return model.ErrUserNotFound
	}

	r.users[user.ID] = user
	return nil
}

// List
func (r *UserRepository) List() ([]*model.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	users := make([]*model.User, 0, len(r.users))
	for _, user := range r.users {
		users = append(users, user)
	}

	return users, nil
}
