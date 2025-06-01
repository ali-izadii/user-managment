package repository

import (
	"context"
	"fmt"
	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"strconv"
	"user-management/internal/config"
)

type PostgresContainer struct {
	Container        *postgres.PostgresContainer
	ConnectionString string
	Host             string
	Port             int
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
