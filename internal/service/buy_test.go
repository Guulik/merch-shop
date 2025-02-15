package service

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"merch/configure"
	"merch/internal/domain/consts"
	"merch/internal/repository/mocks"
	"testing"
	"time"
)

func TestBuy(t *testing.T) {
	type request struct {
		userId int
		item   string
	}
	type data struct {
		customerBalance int
		itemCost        int
	}
	type callErrors struct {
		authorizerError error
		providerError   error
		transferError   error
	}
	tests := []struct {
		name    string
		req     request
		data    data
		errs    callErrors
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			req: request{
				userId: 1,
				item:   "cup",
			},
			data: data{
				customerBalance: 500, //баланса хватает
				itemCost:        20,
			},
			errs:    callErrors{},
			wantErr: assert.NoError,
		},
		{
			name: "user not found",
			req: request{
				userId: 1,
				item:   "cup",
			},
			data: data{
				customerBalance: 0,
				itemCost:        20,
			},
			errs: callErrors{
				providerError: pgx.ErrNoRows,
			},
			wantErr: func(t assert.TestingT, err error, msgAndArgs ...interface{}) bool {
				return assert.Contains(t, err.Error(), consts.UserNotFound)
			},
		},
		{
			name: "provider internal error",
			req: request{
				userId: 1,
				item:   "cup",
			},
			data: data{
				customerBalance: 23,
				itemCost:        20,
			},
			errs: callErrors{
				providerError: errors.New("some strange error"),
			},
			wantErr: func(t assert.TestingT, err error, msgAndArgs ...interface{}) bool {
				return assert.Contains(t, err.Error(), consts.InternalServerError)
			},
		},
		{
			name: "transfer internal error",
			req: request{
				userId: 1,
				item:   "cup",
			},
			data: data{
				customerBalance: 600,
				itemCost:        20,
			},
			errs: callErrors{
				transferError: errors.New("some strange error"),
			},
			wantErr: func(t assert.TestingT, err error, msgAndArgs ...interface{}) bool {
				return assert.Contains(t, err.Error(), consts.InternalServerError)
			},
		},
		{
			name: "error during payment",
			req: request{
				userId: 1,
				item:   "cup",
			},
			data: data{
				customerBalance: 500,
				itemCost:        20,
			},
			errs: callErrors{
				transferError: pgx.ErrNoRows,
			},
			wantErr: func(t assert.TestingT, err error, msgAndArgs ...interface{}) bool {
				return assert.Contains(t, err.Error(), consts.UserNotFound)
			},
		},
		{
			name: "too much:(",
			req: request{
				userId: 1,
				item:   "pink-hoody",
			},
			data: data{
				customerBalance: 228, //баланса НЕ хватает
				itemCost:        500,
			},
			errs: callErrors{},
			wantErr: func(t assert.TestingT, err error, msgAndArgs ...interface{}) bool {
				return assert.Contains(t, err.Error(), consts.NotEnoughMoney)
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
				Return(tt.data.customerBalance, tt.errs.providerError)

			mockTransfer.On("PayForItem", mock.Anything, tt.req.userId, tt.req.item, tt.data.itemCost).
				Return(tt.errs.transferError)

			s := &Service{
				cfg:          cfg,
				authorizer:   mockAuthorizer,
				userProvider: mockProvider,
				coinTransfer: mockTransfer,
			}

			err := s.BuyItem(ctx, tt.req.userId, tt.req.item)
			tt.wantErr(t, err)
		})
	}

}
