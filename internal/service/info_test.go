package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"merch/internal/configure"
	"merch/internal/domain/consts"
	"merch/internal/domain/model"
	"merch/internal/repository/mocks"
)

func TestInfo(t *testing.T) {
	type request struct {
		userId int
	}
	type data struct {
		coins       int
		inventory   map[string]int
		coinHistory model.CoinHistory
	}
	type callErrors struct {
		authorizerError error
		providerError   error
		transferError   error
	}
	tests := []struct {
		name     string
		req      request
		data     data
		errs     callErrors
		wantErr  assert.ErrorAssertionFunc
		wantUser *model.UserInfo
	}{
		{
			name: "success",
			req: request{
				userId: 1,
			},
			data: data{
				coins:     950,
				inventory: map[string]int{"book": 1},
				coinHistory: model.CoinHistory{
					Sent: []model.Sent{model.Sent{
						ToUser: "joji",
						Amount: 50,
					},
					},
				},
			},
			errs:    callErrors{},
			wantErr: assert.NoError,
			wantUser: &model.UserInfo{
				Coins: 950,
				Inventory: []model.Item{model.Item{
					Type:     "book",
					Quantity: 1,
				}},
				CoinHistory: model.CoinHistory{
					Sent: []model.Sent{model.Sent{
						ToUser: "joji",
						Amount: 50,
					},
					},
				},
			},
		},
		{
			name: "provider not found",
			req: request{
				userId: -1,
			},
			data: data{},
			errs: callErrors{
				providerError: pgx.ErrNoRows,
			},
			wantErr: func(t assert.TestingT, err error, msgAndArgs ...interface{}) bool {
				return assert.Contains(t, err.Error(), consts.UserNotFound)
			},
			wantUser: &model.UserInfo{
				Coins:       0,
				Inventory:   nil,
				CoinHistory: model.CoinHistory{},
			},
		},
		{
			name: "provider internal error",
			req: request{
				userId: 1,
			},
			data: data{},
			errs: callErrors{
				providerError: errors.New("some strange error"),
			},
			wantErr: func(t assert.TestingT, err error, msgAndArgs ...interface{}) bool {
				return assert.Contains(t, err.Error(), consts.InternalServerError)
			},
			wantUser: &model.UserInfo{
				Coins:       0,
				Inventory:   nil,
				CoinHistory: model.CoinHistory{},
			},
		},
	}
	ctx := context.Background()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &configure.Config{TokenTTL: time.Minute}
			mockAuthorizer := new(mocks.MockAuthorizer)
			mockProvider := new(mocks.MockUserProvider)
			mockTransfer := new(mocks.MockCoinTransfer)

			mockProvider.On("GetCoins", mock.Anything, tt.req.userId).
				Return(tt.data.coins, tt.errs.providerError)

			mockProvider.On("GetInventory", mock.Anything, tt.req.userId).
				Return(tt.data.inventory, tt.errs.providerError)

			mockProvider.On("GetCoinHistory", mock.Anything, tt.req.userId).
				Return(tt.data.coinHistory, tt.errs.providerError)

			s := &Service{
				cfg:          cfg,
				authorizer:   mockAuthorizer,
				userProvider: mockProvider,
				coinTransfer: mockTransfer,
			}

			user, err := s.GetUserInfo(ctx, tt.req.userId)
			require.Equal(t, tt.wantUser, user)
			tt.wantErr(t, err)
		})
	}

}
