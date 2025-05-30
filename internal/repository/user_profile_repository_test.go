package repository

import (
	"context"
	"testing"
	"time"
	"user-management/internal/config"
	"user-management/internal/database"
	"user-management/internal/models"

	"github.com/google/uuid"
)

func TestUserProfileRepository(t *testing.T) {
	getConfig, err := config.GetConfig("test")
	if err != nil {
		return
	}

	database.InitPostgres(getConfig.Postgres.Uri)
	defer database.Close()

	repo := NewUserProfileRepository()
	ctx := context.Background()

	// Create test user
	id := uuid.New()
	now := time.Now()

	user := models.UserProfile{
		ID:          id,
		Email:       "test.user@example.com",
		FirstName:   "Test",
		LastName:    "User",
		PhoneNumber: "09123456789",
		Bio:         "Test bio",
		IsActive:    true,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	// Insert user
	if err := repo.Create(ctx, user); err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	// Retrieve user
	got, err := repo.GetByID(ctx, id)
	if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}

	if got.Email != user.Email || got.PhoneNumber != user.PhoneNumber {
		t.Errorf("Mismatch in user data. Got %+v, want %+v", got, user)
	}
}
