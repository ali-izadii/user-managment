package repository

import (
	"context"
	"github.com/google/uuid"
	"user-management/internal/models"
)

type UserRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	GetAll(ctx context.Context, limit, offset int) ([]*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	Create(ctx context.Context, user *models.User) error
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type RoleRepository interface {
	Create(ctx context.Context, role *models.Role) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Role, error)
	GetByName(ctx context.Context, name string, organizationID uuid.UUID) (*models.Role, error)
	ListByOrganization(ctx context.Context, organizationID uuid.UUID, limit, offset int) ([]models.Role, error)
	Update(ctx context.Context, role *models.Role) error
	Delete(ctx context.Context, id uuid.UUID) error
	AssignPermissions(ctx context.Context, roleID uuid.UUID, permissionIDs []uuid.UUID) error
	RemovePermissions(ctx context.Context, roleID uuid.UUID, permissionIDs []uuid.UUID) error
	GetRolePermissions(ctx context.Context, roleID uuid.UUID) ([]models.Permission, error)
}

type PermissionRepository interface {
	Create(ctx context.Context, permission *models.Permission) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Permission, error)
	GetByName(ctx context.Context, name string) (*models.Permission, error)
	List(ctx context.Context, limit, offset int) ([]models.Permission, error)
	ListByResource(ctx context.Context, resource string) ([]models.Permission, error)
	Update(ctx context.Context, id uuid.UUID, updates *models.Permission) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type UserRoleRepository interface {
	AssignRole(ctx context.Context, userRole *models.UserRole) error
	RemoveRole(ctx context.Context, userID, roleID, organizationID uuid.UUID) error
	GetUserRoles(ctx context.Context, userID, organizationID uuid.UUID) ([]models.Role, error)
	GetRoleUsers(ctx context.Context, roleID uuid.UUID) ([]models.User, error)
	GetUserPermissions(ctx context.Context, userID, organizationID uuid.UUID) ([]models.Permission, error)
	HasPermission(ctx context.Context, userID uuid.UUID, organizationID uuid.UUID, permission string) (bool, error)
	ListUserOrganizations(ctx context.Context, userID uuid.UUID) ([]models.Organization, error)
}
