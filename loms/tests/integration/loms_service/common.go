package loms_service

import (
	"context"
	"database/sql"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"log"
	"route256.ozon.ru/project/loms/internal/app/server"
	"route256.ozon.ru/project/loms/internal/entity"
	"route256.ozon.ru/project/loms/internal/repository/db"
	"route256.ozon.ru/project/loms/internal/service/loms"
	"route256.ozon.ru/project/loms/pkg/logger"
	"route256.ozon.ru/project/loms/tests/integration/loms_service/testhelper"
)

type StorageTestSuit struct {
	suite.Suite
	storage     loms.Storage
	lomsService server.Service
	pgContainer *testhelper.PostgresContainer
	ctx         context.Context
}

func (s *StorageTestSuit) SetupSuite() {
	s.ctx = context.Background()

	pgContainer, err := testhelper.CreatePostgresContainer(s.ctx)
	if err != nil {
		log.Fatal(err)
	}
	s.pgContainer = pgContainer

	logger.Init()

	dbPool, err := pgxpool.New(s.ctx, s.pgContainer.ConnectionString)

	s.storage = db.NewLomsStorageFromJson(dbPool, dbPool)
	s.lomsService = loms.NewLOMSService(s.storage)
}

func (s *StorageTestSuit) TearDownSuite() {
	if err := s.pgContainer.Terminate(s.ctx); err != nil {
		log.Fatalf("error terminating postgres container: %s", err)
	}
}

func (s *StorageTestSuit) TestCreateOrder() {
	dbConn, err := sql.Open("postgres", s.pgContainer.ConnectionString)
	if err != nil {
		log.Fatal(err)
	}

	MigrateUp(dbConn)
	defer func() {
		MigrateDown(dbConn)
		dbConn.Close()
	}()

	userID := 10
	items := []entity.Item{
		{
			SKU:   1,
			Count: 10,
		},
	}

	orderID, err := s.lomsService.OrderCreate(s.ctx, int64(userID), items)
	require.NoError(s.T(), err)
	require.NotNil(s.T(), orderID)
}

func (s *StorageTestSuit) TestPayOrder() {
	dbConn, err := sql.Open("postgres", s.pgContainer.ConnectionString)
	if err != nil {
		log.Fatal(err)
	}

	MigrateUp(dbConn)
	defer func() {
		MigrateDown(dbConn)
		dbConn.Close()
	}()

	err = s.lomsService.OrderCancel(s.ctx, 1)
	require.NoError(s.T(), err)
}

func (s *StorageTestSuit) TestCancelOrder() {
	dbConn, err := sql.Open("postgres", s.pgContainer.ConnectionString)
	if err != nil {
		log.Fatal(err)
	}

	MigrateUp(dbConn)
	defer func() {
		MigrateDown(dbConn)
		dbConn.Close()
	}()

	err = s.lomsService.OrderPay(s.ctx, 1)
	require.NoError(s.T(), err)
}

func (s *StorageTestSuit) TestOrderInfo() {
	dbConn, err := sql.Open("postgres", s.pgContainer.ConnectionString)
	if err != nil {
		log.Fatal(err)
	}

	MigrateUp(dbConn)
	defer func() {
		MigrateDown(dbConn)
		dbConn.Close()
	}()

	inputOrderID := 1

	expectedInfo := &entity.OrderInfo{
		Status: entity.StatusAwaitingPayment,
		User:   123,
		Items:  []entity.Item{{SKU: 1, Count: 10}},
	}

	info, err := s.lomsService.OrderInfo(s.ctx, entity.OrderID(inputOrderID))
	require.NoError(s.T(), err)
	require.Equal(s.T(), expectedInfo, info)
}

func (s *StorageTestSuit) TestStocksInfo() {
	dbConn, err := sql.Open("postgres", s.pgContainer.ConnectionString)
	if err != nil {
		log.Fatal(err)
	}

	MigrateUp(dbConn)
	defer func() {
		MigrateDown(dbConn)
		dbConn.Close()
	}()

	count, err := s.lomsService.StocksInfo(s.ctx, 1)
	require.NoError(s.T(), err)
	require.Equal(s.T(), entity.Count(90), count)
}
