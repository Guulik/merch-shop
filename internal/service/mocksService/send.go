package mocksService

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MockCoinSender struct {
	mock.Mock
}

func (m *MockCoinSender) SendCoins(ctx context.Context, fromUserId int, toUsername string, coinAmount int) error {
	args := m.Called(ctx, fromUserId, toUsername, coinAmount)
	return args.Error(0)
}
