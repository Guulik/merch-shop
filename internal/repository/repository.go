package repository

import (
	"github.com/jmoiron/sqlx"
)

type Repo struct {
	//log
	dbPool *sqlx.DB
}

func New(pool *sqlx.DB) *Repo {
	return &Repo{
		dbPool: pool,
	}
}
