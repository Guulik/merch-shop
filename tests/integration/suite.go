//go:build integration

package integration

import (
	"context"
	"net/http/httptest"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/stretchr/testify/suite"

	"merch/internal/api"
	"merch/internal/app"
	"merch/internal/configure"
	"merch/internal/repository"
	"merch/internal/service"
)

type Suite struct {
	suite.Suite
	pgContainer *PostgreSQLContainer
	server      *httptest.Server
}

func (s *Suite) SetupSuite() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer ctxCancel()

	pgContainer, err := NewPostgreSQLContainer(ctx)
	s.pgContainer = pgContainer
	s.Require().NoError(err)

	cfg := &configure.Config{TokenTTL: time.Hour}
	pool := configure.NewPostgresPool(ctx, pgContainer.ConnectionString())
	repo := repository.New(pool)
	svc := service.New(cfg, repo, repo, repo)
	handlers := api.New(svc, svc, svc, svc)

	migrationsURL := "file://../../migrations/up"
	m, err := migrate.New(migrationsURL, pgContainer.ConnectionString())
	s.Require().NoError(err)
	err = m.Up()
	s.Require().NoError(err)

	e := app.SetupEcho(handlers)
	s.server = httptest.NewServer(e)
}

func (s *Suite) TearDownSuite() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer ctxCancel()

	s.Require().NoError(s.pgContainer.Terminate(ctx))

	s.server.Close()
}
