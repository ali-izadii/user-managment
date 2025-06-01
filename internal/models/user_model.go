package models

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID            uuid.UUID `db:"id" json:"id"`
	Email         string    `db:"email" json:"email"`
	Password      string    `db:"password" json:"-"`
	FirstName     string    `db:"first_name" json:"first_name"`
	LastName      string    `db:"last_name" json:"last_name"`
	Bio           string    `db:"bio" json:"bio,omitempty"`
	PhoneNumber   string    `db:"phone_number" json:"phone_number,omitempty"`
	EmailVerified bool      `db:"email_verified" json:"email_verified"`
	IsActive      bool      `db:"is_active" json:"is_active"`
	LastLoginAt   time.Time `db:"last_login_at" json:"last_login_at,omitempty"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time `db:"updated_at" json:"updated_at"`
}
