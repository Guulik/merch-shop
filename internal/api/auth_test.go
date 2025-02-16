package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"merch/internal/service/mocksService"
)

func TestAuthHandler(t *testing.T) {
	tests := []struct {
		name          string
		reqBody       string
		expectCode    int
		authToken     string
		authorizerErr error
		wantErr       assert.ErrorAssertionFunc
	}{
		{
			name:          "success - valid credentials",
			reqBody:       `{"username": "test", "password": "secret"}`,
			expectCode:    http.StatusOK,
			authToken:     "valid_token",
			authorizerErr: nil,
			wantErr:       assert.NoError,
		},
		{
			name:       "error - bad request: incorrect format",
			reqBody:    `{"username": 123, "password": "secret"}`,
			expectCode: http.StatusBadRequest,
			wantErr: func(t assert.TestingT, err error, msgAndArgs ...interface{}) bool {
				return assert.Contains(t, err.Error(), "Unmarshal type error")
			},
		},
		{
			name:       "error - bad request: no username",
			reqBody:    `{"password": "secret"}`,
			expectCode: http.StatusBadRequest,
			wantErr: func(t assert.TestingT, err error, msgAndArgs ...interface{}) bool {
				return assert.Contains(t, err.Error(), "Field validation")
			},
		},
		{
			name:       "error - bad request: no password",
			reqBody:    `{"username": "lazzy2wice"}`,
			expectCode: http.StatusBadRequest,
			wantErr: func(t assert.TestingT, err error, msgAndArgs ...interface{}) bool {
				return assert.Contains(t, err.Error(), "Field validation")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/auth",
				strings.NewReader(tt.reqBody))
			req.Header.Set("Content-Type", "application/json")
			res := httptest.NewRecorder()
			ctx := e.NewContext(req, res)

			mockAuth := new(mocksService.MockAuthorizer)
			if tt.authorizerErr != nil || tt.authToken != "" {
				mockAuth.On("Authorize", mock.Anything, "test", mock.Anything).
					Return(tt.authToken, tt.authorizerErr)
			}

			a := &Api{authorizer: mockAuth}
			err := a.AuthHandler(ctx)
			if err != nil {
				e.HTTPErrorHandler(err, ctx)
			}

			require.Equal(t, tt.expectCode, res.Code)
			tt.wantErr(t, err)
		})
	}
}
