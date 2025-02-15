package mocks

import (
	"context"
	"github.com/stretchr/testify/mock"
)

// Mock CoinTransfer
type MockCoinTransfer struct {
	mock.Mock
}

func (m *MockCoinTransfer) TransferCoins(ctx context.Context, fromUserId, toUserId, coinAmount int) error {
	args := m.Called(ctx, fromUserId, toUserId, coinAmount)
	return args.Error(0)
}

func (m *MockCoinTransfer) PayForItem(ctx context.Context, userId int, item string, itemCost int) error {
	args := m.Called(ctx, userId, item, itemCost)
	return args.Error(0)
}
