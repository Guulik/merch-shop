package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo/v4"
	"log"
	"log/slog"
	"merch/configure"
	"merch/internal/api"
	"merch/internal/lib/logger"
	"merch/internal/middleware"
	"merch/internal/repository"
	"merch/internal/service"
	"net/http"
)

type App struct {
	api  *api.Api
	svc  *service.Service
	repo *repository.Repo
	echo *echo.Echo
	pool *pgxpool.Pool
	cfg  *configure.Config
}

func New(ctx context.Context) *App {
	app := &App{}

	app.cfg = configure.MustLoadConfig()
	logger.InitLogger(app.cfg.Env)
	app.echo = echo.New()
	app.pool = configure.NewPostgres(ctx, app.cfg)

	app.repo = repository.New(app.pool)
	app.svc = service.New(app.cfg, app.repo, app.repo, app.repo)
	app.api = api.New(app.svc)

	err := app.cfg.MigrateUp()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatal("failed to migrate: " + err.Error())
	}

	app.echo.POST("/api/auth", app.api.AuthHandler)

	protected := app.echo.Group("/api")
	protected.Use(middleware.AuthMiddleware)

	protected.GET("/buy/:item", app.api.BuyHandler)
	protected.POST("/sendCoin", app.api.SendCoinHandler)
	protected.GET("/info", app.api.InfoHandler)

	return app
}

func (a *App) Run() error {
	slog.Info("server running")

	err := a.echo.Start(fmt.Sprintf(":%d", a.cfg.Port))
	if err != nil {
		return err
	}

	return nil
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}

func (a *App) Stop(ctx context.Context) error {
	slog.Info("stopping server...")

	if err := a.echo.Shutdown(ctx); err != nil {
		slog.Error("failed to shutdown server")
		return err
	}

	a.pool.Close()
	return nil
}
