package repository

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v4"
	"github.com/pashagolub/pgxmock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRepo_TransferCoins(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name       string
		fromUserId int
		toUserId   int
		coinAmount int
		mockError  error
		wantErr    bool
	}{
		{
			name:       "success - coins transferred",
			fromUserId: 1,
			toUserId:   2,
			coinAmount: 50,
			mockError:  nil,
			wantErr:    false,
		},
		{
			name:       "error - not enough coins",
			fromUserId: 1,
			toUserId:   2,
			coinAmount: 100,
			mockError:  pgx.ErrNoRows,
			wantErr:    true,
		},
		{
			name:       "error - database failure on subtract",
			fromUserId: 1,
			toUserId:   2,
			coinAmount: 50,
			mockError:  errors.New("db failure on subtract"),
			wantErr:    true,
		},
		{
			name:       "error - database failure on add",
			fromUserId: 1,
			toUserId:   2,
			coinAmount: 50,
			mockError:  errors.New("db failure on add"),
			wantErr:    true,
		},
		{
			name:       "error - database failure on insert transaction",
			fromUserId: 1,
			toUserId:   2,
			coinAmount: 50,
			mockError:  errors.New("db failure on insert transaction"),
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB, err := pgxmock.NewPool()
			require.NoError(t, err)
			defer mockDB.Close()

			mockDB.ExpectBeginTx(pgx.TxOptions{IsoLevel: pgx.ReadCommitted})

			subtractQuery := `^UPDATE users SET coins = coins - \$1 WHERE id = \$2 AND coins >= \$1`
			addQuery := `^UPDATE users SET coins = coins \+ \$1 WHERE id = \$2`
			insertTransactionQuery := `^INSERT INTO transactions \(from_user, to_user, amount\) VALUES \(\$1, \$2, \$3\)`

			if tt.mockError != nil {
				mockDB.ExpectExec(subtractQuery).
					WithArgs(tt.coinAmount, tt.fromUserId).
					WillReturnError(tt.mockError)
			} else {
				mockDB.ExpectExec(subtractQuery).
					WithArgs(tt.coinAmount, tt.fromUserId).
					WillReturnResult(pgxmock.NewResult("UPDATE", 1))
			}

			if tt.mockError == nil {
				mockDB.ExpectExec(addQuery).
					WithArgs(tt.coinAmount, tt.toUserId).
					WillReturnResult(pgxmock.NewResult("UPDATE", 1))
			}

			if tt.mockError == nil {
				mockDB.ExpectExec(insertTransactionQuery).
					WithArgs(tt.fromUserId, tt.toUserId, tt.coinAmount).
					WillReturnResult(pgxmock.NewResult("INSERT", 1))
			}

			// Expect commit or rollback depending on the error
			if tt.wantErr {
				mockDB.ExpectRollback()
			} else {
				mockDB.ExpectCommit()
			}

			repo := &Repo{dbPool: mockDB}
			err = repo.TransferCoins(ctx, tt.fromUserId, tt.toUserId, tt.coinAmount)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.NoError(t, mockDB.ExpectationsWereMet())
		})
	}
}
