package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/aedobrynin/soa-hw/posts/internal/app"
	"github.com/aedobrynin/soa-hw/posts/internal/logger"
)

func getConfigPath() string {
	var flagConfigPath string

	flag.StringVar(&flagConfigPath, "c", "./.config.yaml", "path to config file")
	flag.Parse()

	return flagConfigPath
}

func initDB(ctx context.Context, config *app.DatabaseConfig) (*pgxpool.Pool, error) {
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

func main() {
	config, err := app.NewConfig(getConfigPath())
	if err != nil {
		log.Fatal(err)
	}

	logger, err := logger.GetLogger(config.App.Debug)
	if err != nil {
		log.Fatal(err)
	}

	pgxPool, err := initDB(context.Background(), &config.Database)
	if err != nil {
		log.Fatal(err)
	}

	a := app.New(pgxPool, logger, config)

	done := make(chan os.Signal, 1)

	signal.Notify(done, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		if err := a.Serve(); err != nil {
			logger.Sugar().Fatal(err)
		}
	}()

	<-done

	a.Shutdown()
}
