package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"

	"merch/internal/domain/model"
)

type MockAuthorizer struct {
	mock.Mock
}

func (m *MockAuthorizer) CheckUserByUsername(ctx context.Context, username string) (*model.UserAuth, error) {
	args := m.Called(ctx, username)
	return args.Get(0).(*model.UserAuth), args.Error(1)
}

func (m *MockAuthorizer) SaveUser(ctx context.Context, username string, password string) (int, error) {
	args := m.Called(ctx, username, password)
	return args.Int(0), args.Error(1)
}
