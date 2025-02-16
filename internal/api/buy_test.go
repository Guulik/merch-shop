package api

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"merch/internal/domain/consts"
	"merch/internal/service/mocksService"
)

func TestBuyHandler(t *testing.T) {
	tests := []struct {
		name       string
		reqBody    string
		item       string
		expectCode int
		buyErr     error
		wantErr    assert.ErrorAssertionFunc
	}{
		{
			name:       "success - valid request",
			reqBody:    `{"item": "book"}`,
			item:       "book",
			expectCode: http.StatusOK,
			buyErr:     nil,
			wantErr:    assert.NoError,
		},
		{
			name:       "success - invalid item",
			reqBody:    `{"item": "glass"}`,
			item:       "glass",
			expectCode: http.StatusBadRequest,
			buyErr:     errors.New(consts.WrongItem),
			wantErr: func(t assert.TestingT, err error, msgAndArgs ...interface{}) bool {
				return assert.Contains(t, err.Error(), consts.WrongItem)
			},
		},
		{
			name:       "error - bad request: incorrect format",
			reqBody:    `{"item": 123}`,
			expectCode: http.StatusBadRequest,
			wantErr: func(t assert.TestingT, err error, msgAndArgs ...interface{}) bool {
				return assert.Contains(t, err.Error(), "Unmarshal type error")
			},
		},
		{
			name:       "error - bad request: empty item",
			reqBody:    `{"item": ""}`,
			expectCode: http.StatusBadRequest,
			wantErr: func(t assert.TestingT, err error, msgAndArgs ...interface{}) bool {
				return assert.Contains(t, err.Error(), "Field validation")
			},
		},
		{
			name:       "error - internal server error",
			reqBody:    `{"item": "book"}`,
			item:       "book",
			expectCode: http.StatusInternalServerError,
			buyErr:     errors.New(consts.InternalServerError),
			wantErr: func(t assert.TestingT, err error, msgAndArgs ...interface{}) bool {
				return assert.Contains(t, err.Error(), consts.InternalServerError)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/buy", strings.NewReader(tt.reqBody))
			req.Header.Set("Content-Type", "application/json")
			res := httptest.NewRecorder()
			ctx := e.NewContext(req, res)

			ctx.Set("user_id", 1)

			mockBuyer := new(mocksService.MockBuyer)
			if tt.buyErr != nil {
				mockBuyer.On("BuyItem", mock.Anything, 1, tt.item).Return(tt.buyErr)
			} else {
				mockBuyer.On("BuyItem", mock.Anything, 1, tt.item).Return(nil)
			}

			a := &Api{buyer: mockBuyer}
			err := a.BuyHandler(ctx)
			if err != nil {
				e.HTTPErrorHandler(err, ctx)
			}

			require.Equal(t, tt.expectCode, res.Code)
			tt.wantErr(t, err)
		})
	}
}

func Test_validateItem(t *testing.T) {
	tests := []struct {
		name    string
		item    string
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "success",
			item:    "book",
			wantErr: assert.NoError,
		},
		{
			name: "error - non existed item",
			item: "полотенце",
			wantErr: func(t assert.TestingT, err error, msgAndArgs ...interface{}) bool {
				return assert.Contains(t, err.Error(), consts.WrongItem)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.wantErr(t, validateItem(tt.item), fmt.Sprintf("validateItem(%v)", tt.item))
		})
	}
}
