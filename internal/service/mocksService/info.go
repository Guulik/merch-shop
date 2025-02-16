package mocksService

import (
	"context"

	"github.com/stretchr/testify/mock"

	"merch/internal/domain/model"
)

type MockInfoProvider struct {
	mock.Mock
}

func (m *MockInfoProvider) GetUserInfo(ctx context.Context, userId int) (*model.UserInfo, error) {
	args := m.Called(ctx, userId)
	return args.Get(0).(*model.UserInfo), args.Error(1)
}
