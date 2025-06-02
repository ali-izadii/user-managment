package repository

import (
	"context"
	"fmt"
	"github.com/docker/go-connections/nat"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
	"user-management/internal/config"
	"user-management/internal/models"
)

type PostgresContainer struct {
	Container        *postgres.PostgresContainer
	ConnectionString string
}

func setupPostgresContainer(ctx context.Context, config config.PostgresConfig) (*PostgresContainer, error) {

	postgresContainer, err := postgres.Run(ctx, "postgres",
		postgres.WithDatabase(config.Database),
		postgres.WithUsername(config.Username),
		postgres.WithPassword(config.Password),
		postgres.BasicWaitStrategies(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to start postgres container: %w", err)
	}

	connectionString, err := postgresContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		return nil, fmt.Errorf("failed to get connection string: %w", err)
	}

	_, err = postgresContainer.Host(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get container host: %w", err)
	}

	_, err = postgresContainer.MappedPort(ctx, nat.Port(strconv.Itoa(config.Port)))
	if err != nil {
		return nil, fmt.Errorf("failed to get container port: %w", err)
	}

	return &PostgresContainer{
		Container:        postgresContainer,
		ConnectionString: connectionString,
	}, nil
}

func runSQLFiles(ctx context.Context, db *pgxpool.Pool, sqlDir string) error {
	files, err := os.ReadDir(sqlDir)
	if err != nil {
		return fmt.Errorf("failed to read SQL directory %s: %w", sqlDir, err)
	}

	var sqlFiles []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".sql") {
			sqlFiles = append(sqlFiles, file.Name())
		}
	}

	sort.Strings(sqlFiles)

	if len(sqlFiles) == 0 {
		fmt.Printf("No .sql files found in %s\n", sqlDir)
		return nil
	}

	for _, filename := range sqlFiles {
		fmt.Printf("Executing %s...\n", filename)

		filePath := filepath.Join(sqlDir, filename)
		sqlContent, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("failed to read SQL file %s: %w", filename, err)
		}

		// Execute SQL
		_, err = db.Exec(ctx, string(sqlContent))
		if err != nil {
			return fmt.Errorf("failed to execute SQL file %s: %w", filename, err)
		}

		fmt.Printf("✓ %s executed successfully\n", filename)
	}

	fmt.Printf("All SQL files executed! (%d files)\n", len(sqlFiles))
	return nil
}

func newUser(id uuid.UUID, email string) *models.User {
	return &models.User{
		ID:            id,
		Email:         email,
		Password:      "pass",
		FirstName:     "Ali",
		LastName:      "Izadi",
		Bio:           models.StringPtr("Toole"),
		PhoneNumber:   models.StringPtr("09170777331"),
		EmailVerified: true,
		IsActive:      true,
		LastLoginAt:   models.TimePtr(time.Now()),
	}
}
