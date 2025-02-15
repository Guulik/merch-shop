package repository

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/pashagolub/pgxmock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRepo_PayForItem(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name      string
		userId    int
		item      string
		itemCost  int
		mockError error
		mockRows  *pgxmock.Rows
		wantErr   bool
	}{
		{
			name:      "success - item purchased",
			userId:    1,
			item:      "book",
			itemCost:  50,
			mockError: nil,
			wantErr:   false,
		},
		{
			name:      "error - not enough coins",
			userId:    1,
			item:      "pink-hoody",
			itemCost:  500,
			mockError: pgx.ErrNoRows,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB, err := pgxmock.NewPool()
			require.NoError(t, err)
			defer mockDB.Close()

			mockDB.ExpectBeginTx(pgx.TxOptions{IsoLevel: pgx.Serializable})

			coinQuery := `^UPDATE users SET coins = coins - \$1 WHERE id = \$2 AND coins >= \$1 `
			itemQuey := `^INSERT INTO inventory \(user_id, item, quantity\) VALUES \(\$1, \$2, \$3\)
						ON CONFLICT \(user_id, item\) DO UPDATE SET quantity = inventory.quantity \+ EXCLUDED.quantity;$`

			if tt.mockError != nil {
				mockDB.ExpectExec(coinQuery).
					WithArgs(tt.itemCost, tt.userId).
					WillReturnError(tt.mockError)
			} else {
				mockDB.ExpectExec(coinQuery).
					WithArgs(tt.itemCost, tt.userId).
					WillReturnResult(pgxmock.NewResult("UPDATE", 1))
			}
			if tt.mockError == nil {
				mockDB.ExpectExec(itemQuey).
					WithArgs(tt.userId, tt.item, 1).
					WillReturnResult(pgxmock.NewResult("UPDATE", 1))
			}

			if tt.wantErr {
				mockDB.ExpectRollback()
			} else {
				mockDB.ExpectCommit()
			}

			repo := &Repo{dbPool: mockDB}
			err = repo.PayForItem(ctx, tt.userId, tt.item, tt.itemCost)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.NoError(t, mockDB.ExpectationsWereMet())
		})
	}
}
