package app

import (
	"awesomeProject/configure"
	"awesomeProject/internal/middleware"
	"context"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"net/http"
)

type App struct {
	//api     *api.Api
	//svc     *service.Service
	//repo *repository.Repo
	echo *echo.Echo
	pool *sqlx.DB
	cfg  *configure.Config
}

func New() *App {
	app := &App{}

	app.cfg = configure.MustLoadConfig()
	app.echo = echo.New()

	app.pool = configure.NewPostgres(app.cfg)

	//app.svc = service.New(app.storage, app.cache)

	//app.api = api.New(app.svc)

	err := app.cfg.MigrateUp()
	if err != nil {
		//TODO: handle error
	}

	//app.echo.POST("/api/auth", api.AuthHandler)

	protected := app.echo.Group("/api")
	protected.Use(middleware.AuthMiddleware)

	//protected.POST("/sendCoin", api.SendCoinHandler)
	//protected.POST("/createBanner", api.CreateBanner)

	return app
}

func (a *App) Run() error {
	fmt.Println("server running")

	err := a.echo.Start(":4444")
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
	fmt.Println("stopping server..." + " op = app.Stop")

	if err := a.echo.Shutdown(ctx); err != nil {
		fmt.Println("failed to shutdown server")
		return err
	}

	if err := a.pool.Close(); err != nil {
		fmt.Println("failed to close connection")
	}
	return nil
}
