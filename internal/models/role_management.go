package models

import (
	"github.com/google/uuid"
	"time"
)

type Organization struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Slug        string    `json:"slug" db:"slug"`
	Description string    `json:"description" db:"description"`
	IsActive    bool      `json:"is_active" db:"is_active"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type Role struct {
	ID             uuid.UUID `json:"id" db:"id"`
	Name           string    `json:"name" db:"name"`
	Description    string    `json:"description" db:"description"`
	OrganizationID uuid.UUID `json:"organization_id" db:"organization_id"`
	IsSystemRole   bool      `json:"is_system_role" db:"is_system_role"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}

type Permission struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Resource    string    `json:"resource" db:"resource"`
	Action      string    `json:"action" db:"action"`
	Description string    `json:"description" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

type RolePermission struct {
	RoleID       uuid.UUID `json:"role_id" db:"role_id"`
	PermissionID uuid.UUID `json:"permission_id" db:"permission_id"`
	GrantedAt    time.Time `json:"granted_at" db:"granted_at"`
}

type UserRole struct {
	UserID         uuid.UUID  `json:"user_id" db:"user_id"`
	RoleID         uuid.UUID  `json:"role_id" db:"role_id"`
	OrganizationID uuid.UUID  `json:"organization_id" db:"organization_id"`
	AssignedAt     time.Time  `json:"assigned_at" db:"assigned_at"`
	AssignedBy     *uuid.UUID `json:"assigned_by" db:"assigned_by"`
}

type UserOrganization struct {
	UserID         uuid.UUID `json:"user_id" db:"user_id"`
	OrganizationID uuid.UUID `json:"organization_id" db:"organization_id"`
	JoinedAt       time.Time `json:"joined_at" db:"joined_at"`
	Status         string    `json:"status" db:"status"` // active, inactive, pending
}
