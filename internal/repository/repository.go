package repository

import (
	"github.com/jackc/pgx/v4/pgxpool"
)

type Repo struct {
	dbPool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Repo {
	return &Repo{
		dbPool: pool,
	}
}
