package service

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"merch/internal/configure"
	"merch/internal/domain/consts"
	"merch/internal/domain/model"
	"merch/internal/repository/mocks"
	"merch/internal/util/hasher"
	"merch/internal/util/jwtManager"
)

func TestAuthorize(t *testing.T) {
	type request struct {
		username string
		password string
	}
	type data struct {
		existedUser *model.UserAuth
		newUserId   int
	}
	type callErrors struct {
		checkError error
		saveError  error
	}
	tests := []struct {
		name    string
		req     request
		data    data
		errs    callErrors
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "success - existed user",
			req: request{
				username: "joja",
				password: "1",
			},
			data: data{
				existedUser: &model.UserAuth{
					Id:       1,
					Username: "joja",
				},
				newUserId: -1,
			},
			errs:    callErrors{},
			wantErr: assert.NoError,
		},
		{
			name: "success - new user",
			req: request{
				username: "newUser",
				password: "newPassword",
			},
			data: data{
				existedUser: nil,
				newUserId:   50,
			},
			errs: callErrors{
				checkError: pgx.ErrNoRows,
			},
			wantErr: assert.NoError,
		},
		{
			name: "error - incorrect password",
			req: request{
				username: "joja",
				password: "wrongPassword",
			},
			data: data{
				existedUser: &model.UserAuth{
					Id:       1,
					Username: "joja",
				},
				newUserId: -1,
			},
			errs: callErrors{},
			wantErr: func(t assert.TestingT, err error, msgAndArgs ...interface{}) bool {
				return assert.Contains(t, err.Error(), consts.WrongPassword)
			},
		},
		{
			name: "error - internal error during checking",
			req: request{
				username: "yoyo",
				password: "ggg",
			},
			data: data{
				existedUser: &model.UserAuth{
					Id:       1,
					Username: "joja",
				},
				newUserId: -1,
			},
			errs: callErrors{
				checkError: errors.New("some internal error"),
			},
			wantErr: func(t assert.TestingT, err error, msgAndArgs ...interface{}) bool {
				return assert.Contains(t, err.Error(), consts.InternalServerError)
			},
		},
		{
			name: "error - internal error during user creation",
			req: request{
				username: "newUser",
				password: "ggg",
			},
			data: data{
				existedUser: nil,
				newUserId:   44,
			},
			errs: callErrors{
				checkError: pgx.ErrNoRows,
				saveError:  fmt.Errorf("some internal error"),
			},
			wantErr: func(t assert.TestingT, err error, msgAndArgs ...interface{}) bool {
				return assert.Contains(t, err.Error(), consts.InternalServerError)
			},
		},
	}
	ctx := context.Background()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalValue, exists := os.LookupEnv("JWT_SECRET")
			if !exists {
				os.Setenv("JWT_SECRET", "lazzy2wice")
			}
			t.Cleanup(func() {
				if exists {
					os.Setenv("JWT_SECRET", originalValue)
				} else {
					os.Unsetenv("JWT_SECRET")
				}
			})

			var (
				generatedToken string
				err            error
			)

			cfg := &configure.Config{TokenTTL: time.Hour}
			mockAuthorizer := new(mocks.MockAuthorizer)
			mockProvider := new(mocks.MockUserProvider)
			mockTransfer := new(mocks.MockCoinTransfer)

			if tt.data.existedUser != nil {
				tt.data.existedUser.PasswordDb, err = hasher.HashPassword("1")
				require.NoError(t, err)
			}

			if tt.data.newUserId == -1 {
				mockAuthorizer.On("CheckUserByUsername", mock.Anything, tt.req.username).
					Return(tt.data.existedUser, tt.errs.checkError)
			} else {
				mockAuthorizer.On("CheckUserByUsername", mock.Anything, tt.req.username).
					Return(&model.UserAuth{}, tt.errs.checkError)
				mockAuthorizer.On("SaveUser", mock.Anything, tt.req.username, mock.Anything).
					Return(tt.data.newUserId, tt.errs.saveError)
			}

			s := &Service{
				cfg:          cfg,
				authorizer:   mockAuthorizer,
				userProvider: mockProvider,
				coinTransfer: mockTransfer,
			}

			if tt.data.newUserId == -1 {
				generatedToken, err = jwtManager.GenerateJWT(tt.data.existedUser.Id, cfg.TokenTTL)
				require.NoError(t, err)
			} else {
				generatedToken, err = jwtManager.GenerateJWT(tt.data.newUserId, cfg.TokenTTL)
				require.NoError(t, err)
			}
			token, err := s.Authorize(ctx, tt.req.username, tt.req.password)
			tt.wantErr(t, err)

			if err == nil {
				secret, err := jwtManager.FetchSecretKey()

				parsedToken, err := jwtManager.ParseToken(generatedToken, secret)
				require.NoError(t, err)
				parsedGeneratedToken := parsedToken.Claims.(jwt.MapClaims)

				parsedReturnedToken, err := jwtManager.ParseToken(token, secret)
				require.NoError(t, err)
				parsedReturnedTokenClaims := parsedReturnedToken.Claims.(jwt.MapClaims)

				require.Equal(t, parsedGeneratedToken["user_id"], parsedReturnedTokenClaims["user_id"])
			}
		})
	}

}
