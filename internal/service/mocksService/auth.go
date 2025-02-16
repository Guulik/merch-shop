package mocksService

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MockAuthorizer struct {
	mock.Mock
}

func (m *MockAuthorizer) Authorize(ctx context.Context, username, password string) (string, error) {
	args := m.Called(ctx, username, password)
	return args.String(0), args.Error(1)
}
