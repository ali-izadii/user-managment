package repository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"testing"
	"user-management/internal/config"
	"user-management/internal/database"
)

//type CalculatorTestSuite struct {
//	suite.Suite
//	StartingNumber int
//}
//
//// this function executes before the test suite begins execution
//func (suite *CalculatorTestSuite) SetupSuite() {
//	// set StartingNumber to one
//	fmt.Println(">>> From SetupSuite")
//	suite.StartingNumber = 1
//}
//
//// this function executes after all tests executed
//func (suite *CalculatorTestSuite) TearDownSuite() {
//	fmt.Println(">>> From TearDownSuite")
//}
//
//// this function executes before each test case
//func (suite *CalculatorTestSuite) SetupTest() {
//	// reset StartingNumber to one
//	fmt.Println("-- From SetupTest")
//	suite.StartingNumber = 1
//}
//
//// this function executes after each test case
//func (suite *CalculatorTestSuite) TearDownTest() {
//	fmt.Println("-- From TearDownTest")
//}
//
//func (suite *CalculatorTestSuite) TestAddOne() {
//	fmt.Println("From TestAddOne")
//	suite.StartingNumber += 1
//	suite.Equal(2, suite.StartingNumber)
//}
//
////func (suite *CalculatorTestSuite) TestSubtractOne() {
////	fmt.Println("From TestSubtractOne")
////	suite.StartingNumber -= 1
////	suite.Equal(0, suite.StartingNumber)
////}
//
//func TestCalculatorTestSuite(t *testing.T) {
//	suite.Run(t, new(CalculatorTestSuite))
//}

type UserRepositoryTestSuite struct {
	suite.Suite
	container   testcontainers.Container
	repo        UserRepository
	db          *pgxpool.Pool
	pgContainer *PostgresContainer
}

func (suite *UserRepositoryTestSuite) SetupSuite() {
	fmt.Println(">>> From SetupSuite")

	ctx := context.Background()
	conf, err := config.GetConfig("test")
	if err != nil {
		return
	}

	pgContainer, err := setupPostgresContainer(ctx, conf.Postgres)
	require.NoError(suite.T(), err)
	suite.pgContainer = pgContainer

	db, err := database.NewDatabaseConnectionString(ctx, suite.pgContainer.ConnectionString, conf.Postgres)
	require.NoError(suite.T(), err)
	suite.db = db

	repo := NewUserRepository(db)
	suite.repo = repo
}

func (suite *UserRepositoryTestSuite) TearDownSuite() {
	fmt.Println(">>> From TearDownSuite")
}

func TestCalculatorTestSuite(t *testing.T) {
	suite.Run(t, new(UserRepositoryTestSuite))
}

func (suite *UserRepositoryTestSuite) TestAddOne() {
	fmt.Println("From TestAddOne")
}
