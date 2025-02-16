package repository

import (
	"context"
	"errors"
	"testing"

	"github.com/jackc/pgx/v4"
	"github.com/pashagolub/pgxmock"
	"github.com/stretchr/testify/require"

	"merch/internal/domain/model"
)

func TestRepo_GetCoins(t *testing.T) {
	tests := []struct {
		name      string
		userId    int
		mockRows  *pgxmock.Rows
		mockError error
		wantCoins int
		wantErr   bool
	}{
		{
			name:      "success",
			userId:    1,
			mockRows:  pgxmock.NewRows([]string{"coins"}).AddRow(100),
			wantCoins: 100,
			wantErr:   false,
		},
		{
			name:      "error - user not found",
			userId:    2,
			mockRows:  pgxmock.NewRows([]string{}),
			wantCoins: -1,
			wantErr:   true,
		},
		{
			name:      "error - database failure",
			userId:    1,
			mockError: errors.New("db connection error"),
			wantCoins: -1,
			wantErr:   true,
		},
	}
	ctx := context.Background()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB, err := pgxmock.NewPool()
			require.NoError(t, err)
			defer mockDB.Close()

			query := `^SELECT coins FROM users WHERE id = \$1`

			if tt.mockError != nil {
				mockDB.ExpectQuery(query).
					WithArgs(tt.userId).
					WillReturnError(tt.mockError)
			} else {
				mockDB.ExpectQuery(query).
					WithArgs(tt.userId).
					WillReturnRows(tt.mockRows)
			}

			repo := New(mockDB)
			coins, err := repo.GetCoins(ctx, tt.userId)

			if tt.wantErr {
				require.Error(t, err)
				require.Equal(t, tt.wantCoins, coins)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.wantCoins, coins)
			}

			require.NoError(t, mockDB.ExpectationsWereMet())
		})
	}
}

func TestRepo_GetCoinsAndInventory(t *testing.T) {
	tests := []struct {
		name      string
		userId    int
		mockRows  *pgxmock.Rows
		mockError error
		wantInv   map[string]int
		wantErr   bool
	}{
		{
			name:     "success - user found with items",
			userId:   1,
			mockRows: pgxmock.NewRows([]string{"item", "quantity"}).AddRow("book", 3).AddRow("pen", 5),
			wantInv:  map[string]int{"book": 3, "pen": 5},
			wantErr:  false,
		},
		{
			name:     "success - user found with no items",
			userId:   2,
			mockRows: pgxmock.NewRows([]string{"item", "quantity"}).AddRow("", 0),
			wantInv:  map[string]int{},
			wantErr:  false,
		},
		{
			name:      "error - user not found",
			userId:    3,
			mockRows:  pgxmock.NewRows([]string{}),
			wantInv:   nil,
			mockError: pgx.ErrNoRows,
			wantErr:   true,
		},
		{
			name:      "error - database failure",
			userId:    1,
			mockError: errors.New("db connection error"),
			wantInv:   nil,
			wantErr:   true,
		},
	}
	ctx := context.Background()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB, err := pgxmock.NewPool()
			require.NoError(t, err)
			defer mockDB.Close()

			query := `^SELECT item, quantity FROM inventory WHERE user_id = \$1`

			if tt.mockError != nil {
				mockDB.ExpectQuery(query).
					WithArgs(tt.userId).
					WillReturnError(tt.mockError)
			} else {
				mockDB.ExpectQuery(query).
					WithArgs(tt.userId).
					WillReturnRows(tt.mockRows)
			}

			repo := New(mockDB)
			inv, err := repo.GetInventory(ctx, tt.userId)

			if tt.wantErr {
				require.Error(t, err)
				require.ErrorIs(t, err, tt.mockError)
				require.Equal(t, tt.wantInv, inv)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.wantInv, inv)
			}

			require.NoError(t, mockDB.ExpectationsWereMet())
		})
	}
}

func TestRepo_GetCoinHistory(t *testing.T) {
	tests := []struct {
		name        string
		userId      int
		mockRows    *pgxmock.Rows
		mockError   error
		wantHistory model.CoinHistory
		wantErr     bool
	}{
		{
			name:   "success - user has both received and sent transactions",
			userId: 1,
			mockRows: pgxmock.NewRows([]string{
				"fromUserId", "fromUsername", "toUserId", "toUsername", "amount",
			}).AddRow(2, "лёха", 1, "саша", 50).AddRow(1, "саша", 3, "вася", 30),
			wantHistory: model.CoinHistory{
				Received: []model.Received{
					{FromUser: "лёха", Amount: 50},
				},
				Sent: []model.Sent{
					{ToUser: "вася", Amount: 30},
				},
			},
			wantErr: false,
		},
		{
			name:   "success - user has only received transactions",
			userId: 2,
			mockRows: pgxmock.NewRows([]string{
				"fromUserId", "fromUsername", "toUserId", "toUsername", "amount",
			}).AddRow(1, "ваня", 2, "саша", 50),
			wantHistory: model.CoinHistory{
				Received: []model.Received{
					{FromUser: "ваня", Amount: 50},
				},
				Sent: nil,
			},
			wantErr: false,
		},
		{
			name:   "success - user has only sent transactions",
			userId: 1,
			mockRows: pgxmock.NewRows([]string{
				"fromUserId", "fromUsername", "toUserId", "toUsername", "amount",
			}).AddRow(1, "саша", 3, "антоха", 30),
			wantHistory: model.CoinHistory{
				Received: nil,
				Sent: []model.Sent{
					{ToUser: "антоха", Amount: 30},
				},
			},
			wantErr: false,
		},
		{
			name:      "error - database failure",
			userId:    1,
			mockError: errors.New("db connection error"),
			wantHistory: model.CoinHistory{
				Received: nil,
				Sent:     nil,
			},
			wantErr: true,
		},
		{
			name:     "error - no transactions found",
			userId:   3,
			mockRows: pgxmock.NewRows([]string{}),
			wantHistory: model.CoinHistory{
				Received: nil,
				Sent:     nil,
			},
			wantErr: false,
		},
	}

	ctx := context.Background()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB, err := pgxmock.NewPool()
			require.NoError(t, err)
			defer mockDB.Close()

			query := `^SELECT t.from_user AS fromUserId, from_user.username AS fromUsername, t.to_user AS toUserId, to_user.username AS toUsername, t.amount
FROM transactions t 
JOIN users from_user ON t.from_user = from_user.id 
JOIN users to_user ON t.to_user = to_user.id 
WHERE t.from_user = \$1 OR t.to_user = \$1`

			if tt.mockError != nil {
				mockDB.ExpectQuery(query).
					WithArgs(tt.userId).
					WillReturnError(tt.mockError)
			} else {
				mockDB.ExpectQuery(query).
					WithArgs(tt.userId).
					WillReturnRows(tt.mockRows)
			}

			repo := New(mockDB)
			history, err := repo.GetCoinHistory(ctx, tt.userId)

			if tt.wantErr {
				require.Error(t, err)
				require.Equal(t, tt.wantHistory, history)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.wantHistory, history)
			}

			require.NoError(t, mockDB.ExpectationsWereMet())
		})
	}
}
