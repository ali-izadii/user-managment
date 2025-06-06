package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"testing"
	"user-management/internal/config"
	"user-management/internal/database"
)

type UserRepositoryTestSuite struct {
	suite.Suite
	repo        UserRepository
	db          *pgxpool.Pool
	pgContainer *PostgresContainer
	ctx         context.Context
}

func (suite *UserRepositoryTestSuite) SetupSuite() {
	suite.ctx = context.Background()
	conf, err := config.GetConfig("test")
	if err != nil {
		return
	}

	pgContainer, err := setupPostgresContainer(suite.ctx, conf.Postgres)
	require.NoError(suite.T(), err)
	suite.pgContainer = pgContainer

	db, err := database.NewDatabaseConnectionString(suite.ctx, suite.pgContainer.ConnectionString, conf.Postgres)
	require.NoError(suite.T(), err)
	suite.db = db

	repo := NewUserRepository(db)
	suite.repo = repo

	err = runSQLFiles(suite.ctx, suite.db, "../database/migrations")
	if err != nil {
		return
	}
}

func (suite *UserRepositoryTestSuite) TearDownSuite() {
	ctx := context.Background()
	if suite.pgContainer != nil && suite.pgContainer.Container != nil {
		err := suite.pgContainer.Container.Terminate(ctx)
		require.NoError(suite.T(), err)
	}
}

func (suite *UserRepositoryTestSuite) SetupTest() {
	ctx := context.Background()
	_, err := suite.db.Exec(ctx, "TRUNCATE TABLE users CASCADE")
	require.NoError(suite.T(), err)
}

func TestCalculatorTestSuite(t *testing.T) {
	suite.Run(t, new(UserRepositoryTestSuite))
}

func (suite *UserRepositoryTestSuite) TestGetUserById() {
	id := uuid.New()

	err := suite.repo.Create(suite.ctx, newUser(id, "fake@email"))
	require.NoError(suite.T(), err)

	user, err := suite.repo.GetByID(suite.ctx, id)
	require.NotNil(suite.T(), user)
	require.Equal(suite.T(), "fake@email", user.Email)
}

func (suite *UserRepositoryTestSuite) TestGetAllUsers() {
	err := suite.repo.Create(suite.ctx, newUser(uuid.New(), "fake1@email"))
	require.NoError(suite.T(), err)
	err = suite.repo.Create(suite.ctx, newUser(uuid.New(), "fake2@email"))
	require.NoError(suite.T(), err)

	users, err := suite.repo.GetAll(suite.ctx, 10, 0)
	require.Equal(suite.T(), 2, len(users))
}

func (suite *UserRepositoryTestSuite) TestGetUserByEmail() {
	err := suite.repo.Create(suite.ctx, newUser(uuid.New(), "fake1@email"))
	require.NoError(suite.T(), err)
	err = suite.repo.Create(suite.ctx, newUser(uuid.New(), "fake2@email"))
	require.NoError(suite.T(), err)

	user, err := suite.repo.GetByEmail(suite.ctx, "fake2@email")
	require.Equal(suite.T(), "fake2@email", user.Email)
}
