package integration

import (
	"github.com/stretchr/testify/suite"
	"route256.ozon.ru/project/cart/tests/integration/cart_service"
	"testing"
)

func TestSmokeSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(cart_service.Suit))
}
