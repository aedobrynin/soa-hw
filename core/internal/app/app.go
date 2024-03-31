package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/jackc/pgx/v4/pgxpool"

	"core/internal/httpadapter"
	"core/internal/repo/userrepo"
	"core/internal/service"
	"core/internal/service/authsvc"
	"core/internal/service/usersvc"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type app struct {
	config *Config

	authSerice  service.Auth
	userService service.User

	httpAdapter httpadapter.Adapter
}

func (a *app) Serve() error {
	done := make(chan os.Signal, 1)

	signal.Notify(done, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		if err := a.httpAdapter.Serve(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err.Error())
		}
	}()

	<-done

	a.Shutdown()

	return nil
}

func (a *app) Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), a.config.App.ShutdownTimeout)
	defer cancel()

	a.httpAdapter.Shutdown(ctx)
}

func New(config *Config) (App, error) {
	pgxPool, err := initDB(context.Background(), &config.Database)
	if err != nil {
		return nil, fmt.Errorf("error on db initialization: %v", err)
	}

	userRepo := userrepo.New(&config.Auth, pgxPool)
	authService := authsvc.New(&config.Auth, userRepo)
	userService := usersvc.New(userRepo)

	a := &app{
		config:      config,
		authSerice:  authService,
		userService: userService,
		httpAdapter: httpadapter.NewAdapter(&config.HTTP, authService, userService),
	}

	return a, nil
}

func initDB(ctx context.Context, config *DatabaseConfig) (*pgxpool.Pool, error) {
	pgxConfig, err := pgxpool.ParseConfig(config.DSN)
	if err != nil {
		return nil, err
	}

	pool, err := pgxpool.ConnectConfig(ctx, pgxConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	// migrations

	m, err := migrate.New(config.MigrationsDir, config.DSN)
	if err != nil {
		return nil, fmt.Errorf("error on migrations creation step: %v", err)
	}

	if err := m.Down(); err != nil && err != migrate.ErrNoChange {
		return nil, fmt.Errorf("error on migrations down step: %v", err)
	}

	if err := m.Up(); err != nil {
		return nil, fmt.Errorf("error on migrations up step: %v", err)
	}

	return pool, nil
}