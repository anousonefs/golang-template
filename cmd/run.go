package cmd

import (
	"context"
	"database/sql"
	"encoding/hex"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	config "github.com/anousoneFS/clean-architecture/config"
	"github.com/anousoneFS/clean-architecture/internal/auth"
	"github.com/anousoneFS/clean-architecture/internal/user"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
)

// Run() is for start program
func Run() error {
	// load config
	cfg, err := config.NewConfig()
	if err != nil {
		return err
	}
	ctx := context.Background()

	// tiggle os Interrupt
	errCh := make(chan error, 1)
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, os.Kill)
	defer cancel()

	// connect to database
	db, err := sql.Open(cfg.DBDriver, cfg.DSNInfo())
	if err != nil {
		return err
	}
	if err := db.Ping(); err != nil {
		return err
	}
	defer db.Close()

	e := newEchoServer(cfg)

	repo := user.New(db)
	userUC := user.NewUserUsecase(repo)
	user.NewHandler(e, userUC)

	key, err := hex.DecodeString(os.Getenv("PASETO_SECRET"))
	if len(key) != 32 {
		return err
	}
	authUC := auth.NewAuthUsecase(userUC, key)
	auth.NewHandler(e, authUC)

	go func() {
		errCh <- e.Start(":" + cfg.Port)
	}()

	// Graceful shutdown.
	select {
	case <-ctx.Done():
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := e.Shutdown(ctx); err != nil {
			return fmt.Errorf("shutdown server failure: %v", err)
		}
		fmt.Println("server shutdown")
	case err := <-errCh:
		return fmt.Errorf("start server error: %v", err)
	}
	return nil
}

// newEchoServer creates a new Echo server.
func newEchoServer(cfg config.Config) *echo.Echo {
	mws := []echo.MiddlewareFunc{
		middleware.LoggerWithConfig(middleware.LoggerConfig{
			Skipper: func(c echo.Context) bool {
				return c.Path() == "/" || c.Path() == "/_healthz"
			},
		}),
		middleware.Recover(),
		middleware.Secure(),
		middleware.CORS(),
	}
	e := echo.New()
	e.Static("/assets", fmt.Sprintf("%s%s", cfg.AssetDir, "/assets"))
	e.HideBanner = true
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(mws...)
	e.GET("/_healthz", func(c echo.Context) error { return c.NoContent(http.StatusOK) })

	return e
}
