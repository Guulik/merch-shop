package mocks

import (
	"context"
	"github.com/stretchr/testify/mock"
	"merch/internal/domain/model"
)

// Mock UserProvider
type MockUserProvider struct {
	mock.Mock
}

func (m *MockUserProvider) GetCoins(ctx context.Context, userId int) (int, error) {
	args := m.Called(ctx, userId)
	return args.Int(0), args.Error(1)
}

func (m *MockUserProvider) GetCoinsAndInventory(ctx context.Context, userId int) (*int, map[string]int, error) {
	args := m.Called(ctx, userId)

	var coins *int
	if args.Get(0) != nil {
		coins = args.Get(0).(*int)
	}

	var inventory map[string]int
	if args.Get(1) != nil {
		inventory = args.Get(1).(map[string]int)
	}

	return coins, inventory, args.Error(2)
}

func (m *MockUserProvider) GetCoinHistory(ctx context.Context, userId int) (model.CoinHistory, error) {
	args := m.Called(ctx, userId)
	return args.Get(0).(model.CoinHistory), args.Error(1)
}
