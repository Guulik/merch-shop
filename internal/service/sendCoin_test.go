package service

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"merch/configure"
	"merch/internal/domain/consts"
	"merch/internal/domain/model"
	"merch/internal/repository/mocks"
	"testing"
	"time"
) // Mock Authorizer

func TestSendCoins(t *testing.T) {
	type request struct {
		userId     int
		toUsername string
		coinAmount int
	}
	type data struct {
		receiverId    int
		senderBalance int
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
				userId:     1,
				toUsername: "gegel",
				coinAmount: 100,
			},
			data: data{
				receiverId:    2,
				senderBalance: 500, //баланса хватает
			},
			errs:    callErrors{},
			wantErr: assert.NoError,
		},
		{
			name: "authorizer not found",
			req: request{
				userId:     1,
				toUsername: "jira",
			},
			data: data{
				senderBalance: 23,
				receiverId:    -1,
			},
			errs: callErrors{
				authorizerError: pgx.ErrNoRows,
			},
			wantErr: func(t assert.TestingT, err error, msgAndArgs ...interface{}) bool {
				return assert.Contains(t, err.Error(), consts.ToUserNotFound)
			},
		},
		{
			name: "authorizer internal error",
			req: request{
				userId:     1,
				toUsername: "gegel",
			},
			data: data{
				senderBalance: 23,
				receiverId:    45,
			},
			errs: callErrors{
				authorizerError: errors.New("some strange error"),
			},
			wantErr: func(t assert.TestingT, err error, msgAndArgs ...interface{}) bool {
				return assert.Contains(t, err.Error(), consts.InternalServerError)
			},
		},
		{
			name: "provider not found",
			req: request{
				userId:     -1,
				toUsername: "jora",
			},
			data: data{
				senderBalance: 0,
				receiverId:    2,
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
				userId:     1,
				toUsername: "gegel",
			},
			data: data{
				senderBalance: 23,
				receiverId:    -1,
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
				userId:     1,
				toUsername: "gegel",
			},
			data: data{
				senderBalance: 23,
				receiverId:    -1,
			},
			errs: callErrors{
				transferError: errors.New("some strange error"),
			},
			wantErr: func(t assert.TestingT, err error, msgAndArgs ...interface{}) bool {
				return assert.Contains(t, err.Error(), consts.InternalServerError)
			},
		},

		{
			name: "not enough money",
			req: request{
				userId:     1,
				toUsername: "KA50",
				coinAmount: 7000,
			},
			data: data{
				receiverId:    2,
				senderBalance: 500,
			},
			errs: callErrors{
				transferError: errors.New(consts.NotEnoughMoney),
			},
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

			mockAuthorizer.On("CheckUserByUsername", mock.Anything, tt.req.toUsername).
				Return(&model.UserAuth{Id: tt.data.receiverId}, tt.errs.authorizerError)

			mockProvider.On("GetCoins", mock.Anything, tt.req.userId).
				Return(tt.data.senderBalance, tt.errs.providerError)

			mockTransfer.On("TransferCoins", mock.Anything, tt.req.userId, tt.data.receiverId, tt.req.coinAmount).
				Return(tt.errs.transferError)

			s := &Service{
				cfg:          cfg,
				authorizer:   mockAuthorizer,
				userProvider: mockProvider,
				coinTransfer: mockTransfer,
			}

			err := s.SendCoins(ctx, tt.req.userId, tt.req.toUsername, tt.req.coinAmount)
			tt.wantErr(t, err)
		})
	}

}
