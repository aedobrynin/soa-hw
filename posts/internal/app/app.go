package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"

	"github.com/aedobrynin/soa-hw/posts/internal/grpcadapter"
	"github.com/aedobrynin/soa-hw/posts/internal/repo/postrepo"
	"github.com/aedobrynin/soa-hw/posts/internal/service"
	"github.com/aedobrynin/soa-hw/posts/internal/service/postsvc"
)

type App struct {
	config *Config

	postService service.Post

	grpcAdapter *grpcadapter.Adapter

	logger *zap.Logger
}

func (a *App) Serve() error {
	done := make(chan os.Signal, 1)

	signal.Notify(done, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		if err := a.grpcAdapter.Run(); err != nil {
			defer func() {
				_ = a.logger.Sync()
			}()
			a.logger.Fatal(err.Error())
		}
	}()

	<-done

	a.Shutdown()

	return nil
}

func (a *App) Shutdown() {
	a.grpcAdapter.Stop()
}

func New(logger *zap.Logger, config *Config) (*App, error) {
	pgxPool, err := initDB(context.Background(), &config.Database)
	if err != nil {
		return nil, fmt.Errorf("error on db initialization: %v", err)
	}

	postRepo := postrepo.New(logger, pgxPool)
	postService := postsvc.New(logger, postRepo)

	a := &App{
		logger:      logger,
		config:      config,
		postService: postService,
		grpcAdapter: grpcadapter.New(logger, postService, &config.GRPC),
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

	if err := m.Up(); err != nil {
		return nil, fmt.Errorf("error on migrations up step: %v", err)
	}

	return pool, nil
}
