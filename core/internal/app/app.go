package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/aedobrynin/soa-hw/core/internal/clients/postsclient"
	"github.com/aedobrynin/soa-hw/core/internal/clients/statisticsclient"
	"github.com/aedobrynin/soa-hw/core/internal/config"
	"github.com/aedobrynin/soa-hw/core/internal/httpadapter"
	"github.com/aedobrynin/soa-hw/core/internal/repo"
	"github.com/aedobrynin/soa-hw/core/internal/repo/statisticsrepo"
	"github.com/aedobrynin/soa-hw/core/internal/repo/userrepo"
	"github.com/aedobrynin/soa-hw/core/internal/service"
	"github.com/aedobrynin/soa-hw/core/internal/service/authsvc"
	"github.com/aedobrynin/soa-hw/core/internal/service/statisticssvc"
	"github.com/aedobrynin/soa-hw/core/internal/service/usersvc"
)

type app struct {
	config      *config.Config
	httpAdapter httpadapter.Adapter

	userRepo       repo.User
	statisticsRepo repo.Statistics

	userService       service.User
	authService       service.Auth
	statisticsService service.Statistics
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
	a.statisticsRepo.Stop(ctx)
}

func New(config *config.Config) (App, error) {
	pgxPool, err := initDB(context.Background(), &config.Database)
	if err != nil {
		return nil, fmt.Errorf("error on db initialization: %v", err)
	}

	userRepo := userrepo.New(&config.Auth, pgxPool)
	authService := authsvc.New(&config.Auth, userRepo)
	userService := usersvc.New(userRepo)

	postsClient, err := postsclient.New(context.Background(), &config.Posts)
	if err != nil {
		return nil, fmt.Errorf("error on posts client initialization: %v", err)
	}

	statisticsRepo := statisticsrepo.New(&config.Kafka)
	statisticsService := statisticssvc.New(statisticsRepo)
	statisticsClient, err := statisticsclient.New(context.Background(), &config.Statistics)
	if err != nil {
		return nil, fmt.Errorf("error on statistics client initialization: %v", err)
	}

	a := &app{
		config:            config,
		userRepo:          userRepo,
		statisticsRepo:    statisticsRepo,
		authService:       authService,
		userService:       userService,
		statisticsService: statisticsService,
		httpAdapter: httpadapter.NewAdapter(
			&config.HTTP,
			authService,
			userService,
			statisticsService,
			postsClient,
			statisticsClient,
		),
	}

	return a, nil
}

func initDB(ctx context.Context, config *config.DatabaseConfig) (*pgxpool.Pool, error) {
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

	if err := m.Up(); err != migrate.ErrNoChange {
		return nil, fmt.Errorf("error on migrations up step: %v", err)
	}

	return pool, nil
}
