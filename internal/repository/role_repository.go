package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
	"user-management/internal/models"
)

type roleRepository struct {
	db *pgxpool.Pool
}

func NewRoleRepository(db *pgxpool.Pool) RoleRepository {
	return &roleRepository{db: db}
}

func (r *roleRepository) Create(ctx context.Context, role *models.Role) error {
	if role.ID == uuid.Nil {
		return fmt.Errorf("role ID is required")
	}
	if role.Name == "" {
		return fmt.Errorf("role name is required")
	}
	if role.OrganizationID == uuid.Nil {
		return fmt.Errorf("organization ID is required")
	}

	if role.CreatedAt.IsZero() {
		role.CreatedAt = time.Now()
	}
	if role.UpdatedAt.IsZero() {
		role.UpdatedAt = time.Now()
	}

	query := `
		INSERT INTO roles (id, name, description, organization_id, is_system_role, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := r.db.Exec(ctx, query,
		role.ID,
		role.Name,
		role.Description,
		role.OrganizationID,
		role.IsSystemRole,
		role.CreatedAt,
		role.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create role: %w", err)
	}

	return nil
}

func (r *roleRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Role, error) {
	role := &models.Role{}

	query := `
		SELECT id, name, description, organization_id, is_system_role, created_at, updated_at
		FROM roles
		WHERE id = $1
	`

	err := r.db.QueryRow(ctx, query, id).Scan(
		&role.ID,
		&role.Name,
		&role.Description,
		&role.OrganizationID,
		&role.IsSystemRole,
		&role.CreatedAt,
		&role.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("role not found")
		}
		return nil, fmt.Errorf("failed to get role: %w", err)
	}

	return role, nil
}

func (r *roleRepository) GetByName(ctx context.Context, name string, organizationID uuid.UUID) (*models.Role, error) {
	role := &models.Role{}

	query := `
		SELECT id, name, description, organization_id, is_system_role, created_at, updated_at
		FROM roles
		WHERE name = $1 AND organization_id = $2
	`

	err := r.db.QueryRow(ctx, query, name, organizationID).Scan(
		&role.ID,
		&role.Name,
		&role.Description,
		&role.OrganizationID,
		&role.IsSystemRole,
		&role.CreatedAt,
		&role.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("role not found")
		}
		return nil, fmt.Errorf("failed to get role: %w", err)
	}

	return role, nil
}

func (r *roleRepository) ListByOrganization(ctx context.Context, organizationID uuid.UUID, limit, offset int) ([]models.Role, error) {
	query := `
		SELECT r.id, r.name, r.description, r.organization_id, r.is_system_role, 
		       r.created_at, r.updated_at
		FROM roles r
		WHERE r.organization_id = $1
		ORDER BY r.name
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(ctx, query, organizationID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list roles: %w", err)
	}
	defer rows.Close()

	var roles []models.Role
	for rows.Next() {
		var role models.Role
		err := rows.Scan(
			&role.ID,
			&role.Name,
			&role.Description,
			&role.OrganizationID,
			&role.IsSystemRole,
			&role.CreatedAt,
			&role.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan role: %w", err)
		}
		roles = append(roles, role)
	}

	return roles, nil
}

func (r *roleRepository) Update(ctx context.Context, role *models.Role) error {
	role.UpdatedAt = time.Now()

	query := `
		UPDATE roles 
		SET name = $2, description = $3, updated_at = $4
		WHERE id = $1
	`

	result, err := r.db.Exec(ctx, query,
		role.ID,
		role.Name,
		role.Description,
		role.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to update role: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("role not found")
	}

	return nil
}

func (r *roleRepository) Delete(ctx context.Context, id uuid.UUID) error {
	var isSystemRole bool
	err := r.db.QueryRow(ctx, "SELECT is_system_role FROM roles WHERE id = $1", id).Scan(&isSystemRole)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("role not found")
		}
		return fmt.Errorf("failed to check role: %w", err)
	}

	if isSystemRole {
		return fmt.Errorf("cannot delete system role")
	}

	query := "DELETE FROM roles WHERE id = $1"
	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete role: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("role not found")
	}

	return nil
}

func (r *roleRepository) AssignPermissions(ctx context.Context, roleID uuid.UUID, permissionIDs []uuid.UUID) error {
	if len(permissionIDs) == 0 {
		return nil
	}

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	_, err = tx.Exec(ctx, "DELETE FROM role_permissions WHERE role_id = $1", roleID)
	if err != nil {
		_ = tx.Rollback(ctx)
		return fmt.Errorf("failed to remove existing permissions: %w", err)
	}

	for _, permissionID := range permissionIDs {
		_, err = tx.Exec(ctx,
			"INSERT INTO role_permissions (role_id, permission_id, granted_at) VALUES ($1, $2, $3)",
			roleID, permissionID, time.Now())
		if err != nil {
			_ = tx.Rollback(ctx)
			return fmt.Errorf("failed to assign permission: %w", err)
		}
	}

	return tx.Commit(ctx)
}

func (r *roleRepository) RemovePermissions(ctx context.Context, roleID uuid.UUID, permissionIDs []uuid.UUID) error {
	if len(permissionIDs) == 0 {
		return nil
	}

	query := "DELETE FROM role_permissions WHERE role_id = $1 AND permission_id = ANY($2)"
	_, err := r.db.Exec(ctx, query, roleID, permissionIDs)
	if err != nil {
		return fmt.Errorf("failed to remove permissions: %w", err)
	}

	return nil
}

func (r *roleRepository) GetRolePermissions(ctx context.Context, roleID uuid.UUID) ([]models.Permission, error) {
	query := `
		SELECT p.id, p.name, p.resource, p.action, p.description, p.created_at
		FROM permissions p
		INNER JOIN role_permissions rp ON p.id = rp.permission_id
		WHERE rp.role_id = $1
		ORDER BY p.resource, p.action
	`

	rows, err := r.db.Query(ctx, query, roleID)
	if err != nil {
		return nil, fmt.Errorf("failed to get role permissions: %w", err)
	}
	defer rows.Close()

	var permissions []models.Permission
	for rows.Next() {
		var permission models.Permission
		err := rows.Scan(
			&permission.ID,
			&permission.Name,
			&permission.Resource,
			&permission.Action,
			&permission.Description,
			&permission.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan permission: %w", err)
		}
		permissions = append(permissions, permission)
	}

	return permissions, nil
}
