package repository

import (
	"context"
	"fmt"
	"user-management/internal/database"
	"user-management/internal/models"

	"github.com/google/uuid"
)

type UserProfileRepository struct{}

func NewUserProfileRepository() *UserProfileRepository {
	return &UserProfileRepository{}
}

func (r *UserProfileRepository) Create(ctx context.Context, user models.UserProfile) error {
	query := `
		INSERT INTO user_profiles (
			id, email, first_name, last_name, phone_number, bio, is_active, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9
		)
	`
	_, err := database.Conn.Exec(ctx, query,
		user.ID,
		user.Email,
		user.FirstName,
		user.LastName,
		user.PhoneNumber,
		user.Bio,
		user.IsActive,
		user.CreatedAt,
		user.UpdatedAt,
	)

	return err
}

func (r *UserProfileRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.UserProfile, error) {
	query := `
		SELECT id, email, first_name, last_name, phone_number, bio, is_active, created_at, updated_at
		FROM user_profiles WHERE id = $1
	`

	row := database.Conn.QueryRow(ctx, query, id)

	var user models.UserProfile
	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.PhoneNumber,
		&user.Bio,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	return &user, nil
}
