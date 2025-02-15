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
}

func New(ctx context.Context, cfg *configure.Config) *App {
	app := &App{}

	logger.InitLogger(cfg.Env)

	app.pool = configure.NewPostgresPool(ctx, cfg.ConnectionString())

	app.repo = repository.New(app.pool)
	app.svc = service.New(cfg, app.repo, app.repo, app.repo)
	app.api = api.New(app.svc)

	app.echo = SetupEcho(app.api)

	err := cfg.MigrateUp()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatal("failed to migrate: " + err.Error())
	}

	return app
}

func SetupEcho(api *api.Api) *echo.Echo {
	e := echo.New()
	e.POST("/api/auth", api.AuthHandler)

	protected := e.Group("/api")
	protected.Use(middleware.AuthMiddleware)

	protected.GET("/buy/:item", api.BuyHandler)
	protected.POST("/sendCoin", api.SendCoinHandler)
	protected.GET("/info", api.InfoHandler)

	return e
}

func (a *App) Run(port int) error {
	slog.Info("server running")

	err := a.echo.Start(fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}

	return nil
}

func (a *App) MustRun(port int) {
	if err := a.Run(port); err != nil && !errors.Is(err, http.ErrServerClosed) {
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
