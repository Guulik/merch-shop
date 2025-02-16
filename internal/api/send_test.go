package api

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"merch/internal/domain/consts"
	"merch/internal/service/mocksService"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestSendCoinHandler(t *testing.T) {
	tests := []struct {
		name        string
		reqBody     string
		tokenUserId int
		mockError   error
		expectCode  int
		wantErr     assert.ErrorAssertionFunc
	}{
		{
			name:        "success - coins sent",
			reqBody:     `{"toUser": "recipient", "amount": 100}`,
			tokenUserId: 1,
			expectCode:  http.StatusOK,
			wantErr:     assert.NoError,
		},
		{
			name:        "error - bad request: invalid JSON",
			reqBody:     `{"toUser": 123, "amount": "invalid"}`,
			tokenUserId: 1,
			expectCode:  http.StatusBadRequest,
			wantErr: func(t assert.TestingT, err error, msgAndArgs ...interface{}) bool {
				return assert.Contains(t, err.Error(), "Unmarshal type error")
			},
		},
		{
			name:        "error - bad request: missing user",
			reqBody:     `{"amount": 100}`,
			tokenUserId: 1,
			expectCode:  http.StatusBadRequest,
			wantErr: func(t assert.TestingT, err error, msgAndArgs ...interface{}) bool {
				return assert.Contains(t, err.Error(), "Field validation")
			},
		},
		{
			name:        "error - bad request: missing amount",
			reqBody:     `{"toUser": "илюха"}`,
			tokenUserId: 1,
			expectCode:  http.StatusBadRequest,
			wantErr: func(t assert.TestingT, err error, msgAndArgs ...interface{}) bool {
				return assert.Contains(t, err.Error(), "Field validation")
			},
		},
		{
			name:        "error - internal server error",
			reqBody:     `{"toUser": "recipient", "amount": 100}`,
			tokenUserId: 1,
			mockError:   errors.New(consts.InternalServerError),
			expectCode:  http.StatusInternalServerError,
			wantErr: func(t assert.TestingT, err error, msgAndArgs ...interface{}) bool {
				return assert.Contains(t, err.Error(), consts.InternalServerError)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/send-coin", strings.NewReader(tt.reqBody))
			req.Header.Set("Content-Type", "application/json")
			res := httptest.NewRecorder()
			ctx := e.NewContext(req, res)

			ctx.Set("user_id", tt.tokenUserId)

			mockCoinSender := new(mocksService.MockCoinSender)
			if tt.mockError != nil {
				mockCoinSender.On("SendCoins", mock.Anything, tt.tokenUserId, mock.Anything, mock.Anything).
					Return(tt.mockError)
			} else {
				mockCoinSender.On("SendCoins", mock.Anything, tt.tokenUserId, mock.Anything, mock.Anything).
					Return(nil)
			}

			api := &Api{coinSender: mockCoinSender}
			err := api.SendCoinHandler(ctx)
			if err != nil {
				e.HTTPErrorHandler(err, ctx)
			}

			require.Equal(t, tt.expectCode, res.Code)
			tt.wantErr(t, err)
		})
	}
}
