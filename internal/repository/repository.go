package repository

import (
	"github.com/jmoiron/sqlx"
)

type Repo struct {
	dbPool *sqlx.DB
}

func New(pool *sqlx.DB) *Repo {
	return &Repo{
		dbPool: pool,
	}
}
