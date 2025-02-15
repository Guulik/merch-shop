package mocksService

import (
	"context"
	"github.com/stretchr/testify/mock"
)

type MockBuyer struct {
	mock.Mock
}

func (m *MockBuyer) BuyItem(ctx context.Context, userId int, item string) error {
	args := m.Called(ctx, userId, item)
	return args.Error(0)
}
