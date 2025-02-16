package repository

import (
	"context"

	"github.com/jackc/pgx/v4"
)

type Repo struct {
	//dbPool *pgxpool.Pool
	dbPool PgxPooler
}

func New(pool PgxPooler) *Repo {
	return &Repo{
		dbPool: pool,
	}
}

// PgxPooler - это интерфейс для генерации моков.
// Этот интерфейс реализуется в pgxpool.Pool
// По сути интерфейс "подогнан" под используемые методы из pgxpool.Pool.
type PgxPooler interface {
	Begin(context.Context) (pgx.Tx, error)
	BeginTx(context.Context, pgx.TxOptions) (pgx.Tx, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Close()
}
