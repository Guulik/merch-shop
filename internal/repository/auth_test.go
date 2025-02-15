package repository

import (
	"context"
	"errors"
	"github.com/pashagolub/pgxmock"
	"github.com/stretchr/testify/require"
	"merch/internal/domain/model"
	"testing"
)

func TestRepo_CheckUserByUsername(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name      string
		username  string
		mockRows  *pgxmock.Rows
		mockError error
		wantUser  *model.UserAuth
		wantErr   bool
	}{
		{
			name:     "success - user exists",
			username: "testUser",
			mockRows: pgxmock.NewRows([]string{"id", "username", "password_hash"}).
				AddRow(1, "testUser", "hashedPassword"),
			wantUser: &model.UserAuth{
				Id:         1,
				Username:   "testUser",
				PasswordDb: "hashedPassword",
			},
			wantErr: false,
		},
		{
			name:     "error - user not found",
			username: "unknownUser",
			mockRows: pgxmock.NewRows([]string{}),
			wantUser: nil,
			wantErr:  true,
		},
		{
			name:      "error - database failure",
			username:  "testUser",
			mockError: errors.New("db connection error"),
			wantUser:  nil,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB, err := pgxmock.NewPool()
			require.NoError(t, err)
			defer mockDB.Close()

			if tt.mockError != nil {
				mockDB.ExpectQuery("SELECT id, username, password_hash FROM users").
					WithArgs(tt.username).
					WillReturnError(tt.mockError)
			} else {
				mockDB.ExpectQuery("SELECT id, username, password_hash FROM users").
					WithArgs(tt.username).
					WillReturnRows(tt.mockRows)
			}

			repo := New(mockDB)
			user, err := repo.CheckUserByUsername(ctx, tt.username)

			if tt.wantErr {
				require.Error(t, err)
				require.Nil(t, user)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.wantUser, user)
			}

			require.NoError(t, mockDB.ExpectationsWereMet())
		})
	}
}

func TestRepo_SaveUser(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name       string
		username   string
		password   string
		mockResult *pgxmock.Rows
		mockError  error
		wantUserID int
		wantErr    bool
	}{
		{
			name:       "success - user saved",
			username:   "usr",
			password:   "hashedPassword",
			mockResult: pgxmock.NewRows([]string{"id"}).AddRow(1),
			wantUserID: 1,
			wantErr:    false,
		},
		{
			name:       "error - database failure",
			username:   "user",
			password:   "hashedPassword",
			mockError:  errors.New("db insert error"),
			wantUserID: 0,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB, err := pgxmock.NewPool()
			require.NoError(t, err)
			defer mockDB.Close()

			query := "INSERT INTO users \\(username, password_hash\\) VALUES \\(\\$1, \\$2\\) RETURNING id;"

			mockQuery := mockDB.ExpectQuery(query).
				WithArgs(tt.username, tt.password)

			if tt.mockError != nil {
				mockQuery.WillReturnError(tt.mockError)
			} else {
				mockQuery.WillReturnRows(tt.mockResult)
			}

			repo := New(mockDB)
			userID, err := repo.SaveUser(ctx, tt.username, tt.password)

			if tt.wantErr {
				require.Error(t, err)
				require.Zero(t, userID)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.wantUserID, userID)
			}

			require.NoError(t, mockDB.ExpectationsWereMet())
		})
	}
}
