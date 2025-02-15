package api

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"merch/internal/domain/consts"
	"merch/internal/domain/model"
	"merch/internal/lib/wrapper"
	"merch/internal/service/mocksService"
)

func TestInfoHandler(t *testing.T) {
	tests := []struct {
		name         string
		userID       int
		mockResponse *model.UserInfo
		mockError    error
		expectedCode int
		expectedBody string
		wantErr      assert.ErrorAssertionFunc
	}{
		{
			name:   "success - valid user",
			userID: 1,
			mockResponse: &model.UserInfo{
				Coins: 100,
				Inventory: []model.Item{
					{
						Type:     "book",
						Quantity: 1,
					},
				},
			},
			expectedCode: http.StatusOK,
			expectedBody: `{
							"coins": 100,
							"inventory": [
								{
									"type": "book",
									"quantity": 1
								}
							],
							"coinHistory": {
								"received": null,
								"sent": null
									}
							}`,
			wantErr: assert.NoError,
		},
		{
			name:         "error - user not found",
			userID:       3,
			mockError:    wrapper.WrapHTTPError(errors.New(consts.UserNotFound), http.StatusBadRequest, consts.UserNotFound),
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"message":"user not found"}`,
			wantErr: func(t assert.TestingT, err error, msgAndArgs ...interface{}) bool {
				return assert.Contains(t, err.Error(), consts.UserNotFound)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/info", nil)
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)
			ctx.Set("user_id", tt.userID)

			mockProvider := new(mocksService.MockInfoProvider)
			mockProvider.On("GetUserInfo", mock.Anything, tt.userID).
				Return(tt.mockResponse, tt.mockError)

			a := &Api{infoProvider: mockProvider}
			err := a.InfoHandler(ctx)
			if err != nil {
				e.HTTPErrorHandler(err, ctx)
			}

			require.Equal(t, tt.expectedCode, rec.Code)
			assert.JSONEq(t, tt.expectedBody, rec.Body.String())
			tt.wantErr(t, err)
		})
	}
}
