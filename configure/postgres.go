package configure

import (
	"context"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v4/pgxpool"
	"log/slog"
	"net/url"
)

func NewPostgres(ctx context.Context, cfg *Config) *pgxpool.Pool {
	pool, err := pgxpool.Connect(ctx, cfg.connectionString())
	if err != nil {
		panic("no connection to database")
	}
	return pool
}

func (c *Config) MigrateUp(url ...string) error {
	var sourceURL string
	if url == nil {
		sourceURL = "file://migrations/up"
	} else {
		sourceURL = url[0]
	}
	slog.Info(c.connectionString())
	m, err := migrate.New(sourceURL, c.connectionString())
	if err != nil {
		return err
	}
	if err = m.Up(); err != nil {
		return err
	}

	return nil
}

func (c *Config) connectionString() string {
	u := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(c.Postgres.User, c.Postgres.Password),
		Host:   fmt.Sprintf("%s:%d", c.Postgres.Host, c.Postgres.SQLPort),
		Path:   c.Postgres.DBName,
	}

	q := u.Query()
	q.Set("sslmode", "disable")

	u.RawQuery = q.Encode()

	return u.String()
}
