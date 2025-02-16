package configure

import (
	"context"
	"fmt"
	"log/slog"
	"net/url"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v4/pgxpool"
)

func NewPostgresPool(ctx context.Context, connectionString string) *pgxpool.Pool {
	pool, err := pgxpool.Connect(ctx, connectionString)
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
	slog.Info(c.ConnectionString())
	m, err := migrate.New(sourceURL, c.ConnectionString())
	if err != nil {
		return err
	}
	if err = m.Up(); err != nil {
		return err
	}

	return nil
}

func (c *Config) ConnectionString() string {
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
