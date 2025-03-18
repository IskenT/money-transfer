package memory_test

import (
	"testing"

	"github.com/IskenT/money-transfer/internal/domain/model"
	"github.com/IskenT/money-transfer/internal/infra/repository/memory"
)

func TestUserRepository(t *testing.T) {
	userRepo := memory.NewUserRepository()

	t.Run("GetByID existing user", func(t *testing.T) {
		user, err := userRepo.GetByID("1")
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if user.ID != "1" || user.Name != "Mark" {
			t.Errorf("Expected user Mark with ID 1, got %s with ID %s", user.Name, user.ID)
		}
	})

	t.Run("GetByID non-existing user", func(t *testing.T) {
		_, err := userRepo.GetByID("999")
		if err != model.ErrUserNotFound {
			t.Errorf("Expected error %v, got %v", model.ErrUserNotFound, err)
		}
	})

	t.Run("Update existing user", func(t *testing.T) {
		user, _ := userRepo.GetByID("1")
		originalBalance := user.Balance

		user.Balance += 1000
		err := userRepo.Update(user)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		updatedUser, _ := userRepo.GetByID("1")
		if updatedUser.Balance != originalBalance+1000 {
			t.Errorf("Expected balance %d, got %d", originalBalance+1000, updatedUser.Balance)
		}
	})

	t.Run("Update non-existing user", func(t *testing.T) {
		user := &model.User{ID: "999", Name: "NonExistent", Balance: 1000}
		err := userRepo.Update(user)
		if err != model.ErrUserNotFound {
			t.Errorf("Expected error %v, got %v", model.ErrUserNotFound, err)
		}
	})

	t.Run("List all users", func(t *testing.T) {
		users, err := userRepo.List()
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if len(users) < 3 {
			t.Errorf("Expected at least 3 users, got %d", len(users))
		}

		userMap := make(map[string]*model.User)
		for _, u := range users {
			userMap[u.ID] = u
		}

		if _, exists := userMap["1"]; !exists {
			t.Errorf("Expected user with ID 1 to exist")
		}
		if _, exists := userMap["2"]; !exists {
			t.Errorf("Expected user with ID 2 to exist")
		}
		if _, exists := userMap["3"]; !exists {
			t.Errorf("Expected user with ID 3 to exist")
		}
	})
}
