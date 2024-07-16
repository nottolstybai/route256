package integration

import (
	"github.com/stretchr/testify/suite"
	"route256.ozon.ru/project/loms/tests/integration/loms_service"
	"testing"
)

func TestLomsStorageTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(loms_service.StorageTestSuit))
}
