package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
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
	pool *sqlx.DB
	cfg  *configure.Config
}

func New() *App {
	app := &App{}

	app.cfg = configure.MustLoadConfig()
	logger.InitLogger(app.cfg.Env)
	app.echo = echo.New()
	app.pool = configure.NewPostgres(app.cfg)

	app.repo = repository.New(app.pool)
	app.svc = service.New(app.cfg, app.repo, app.repo, app.repo)
	app.api = api.New(app.svc)

	err := app.cfg.MigrateUp()
	if err != nil {
		//TODO: handle error
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
	fmt.Println("server running")

	err := a.echo.Start(fmt.Sprintf(":%d", a.cfg.Port))
	if err != nil {
		return err
	}

	return nil
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		panic(err)
	}
}

func (a *App) Stop(ctx context.Context) error {
	fmt.Println("stopping server...")

	if err := a.echo.Shutdown(ctx); err != nil {
		fmt.Println("failed to shutdown server")
		return err
	}

	if err := a.pool.Close(); err != nil {
		fmt.Println("failed to close connection")
	}
	return nil
}
