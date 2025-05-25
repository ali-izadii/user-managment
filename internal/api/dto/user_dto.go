package dto

import (
	"github.com/google/uuid"
	"time"
)

type UpdateProfileRequest struct {
	FirstName   string `json:"first_name" validate:"required,min=2,max=50" example:"Ali"`
	LastName    string `json:"last_name" validate:"required,min=2,max=50" example:"Izadi"`
	PhoneNumber string `json:"phone_number" validate:"omitempty,e164" example:"+1234567890"`
	Bio         string `json:"bio" validate:"omitempty,max=500"`
	Timezone    string `json:"timezone" validate:"omitempty"`
}

type UserProfileDTO struct {
	ID          uuid.UUID `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Email       string    `json:"email" example:"ali.ali@example.com"`
	FirstName   string    `json:"first_name" example:"Ali"`
	LastName    string    `json:"last_name" example:"Izadi"`
	PhoneNumber *string   `json:"phone_number,omitempty"`
	Bio         *string   `json:"bio,omitempty"`
	Timezone    *string   `json:"timezone,omitempty"`
	IsActive    bool      `json:"is_active" example:"true"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ErrorResponse struct {
	Status  string                 `json:"status" example:"error"`
	Message string                 `json:"message" example:"Validation failed"`
	Errors  map[string]interface{} `json:"errors,omitempty"`
}
