package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ClickHouse/clickhouse-go/v2"
	_ "github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/clickhouse"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"go.uber.org/zap"

	"github.com/aedobrynin/soa-hw/statistics/internal/grpcadapter"
	"github.com/aedobrynin/soa-hw/statistics/internal/repo/statisticsrepo"
	"github.com/aedobrynin/soa-hw/statistics/internal/service"
	"github.com/aedobrynin/soa-hw/statistics/internal/service/statisticssvc"
)

type AppImpl struct {
	config *Config

	statisticsService service.Statistics

	grpcAdapter *grpcadapter.Adapter

	logger *zap.Logger
}

func (a *AppImpl) Serve() error {
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

func (a *AppImpl) Shutdown() {
	a.grpcAdapter.Stop()
}

func New(logger *zap.Logger, config *Config) (App, error) {
	err := applyMigrations(&config.Database)
	if err != nil {
		return nil, fmt.Errorf("error on applyMigrations: %v", err)
	}

	conn, err := connectToDb(context.Background(), &config.Database)
	if err != nil {
		return nil, fmt.Errorf("error on connection to ClickHouse: %v", err)
	}

	statisticsRepo := statisticsrepo.New(logger, conn)
	statisticsService := statisticssvc.New(logger, statisticsRepo)

	a := &AppImpl{
		logger:            logger,
		config:            config,
		statisticsService: statisticsService,
		grpcAdapter:       grpcadapter.New(logger, statisticsService, &config.GRPC),
	}

	return a, nil
}

func connectToDb(ctx context.Context, config *DatabaseConfig) (clickhouse.Conn, error) {
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{config.Addr},
		Auth: clickhouse.Auth{
			Database: config.Database,
			Username: config.User,
			Password: config.Password,
		},
		ClientInfo: clickhouse.ClientInfo{
			Products: []struct {
				Name    string
				Version string
			}{
				{Name: "statistics-service-go-client", Version: "0.1"},
			},
		},
		Debugf: func(format string, v ...interface{}) {
			log.Printf(format, v)
		},
	})

	if err != nil {
		return nil, err
	}

	if err := conn.Ping(ctx); err != nil {
		if exception, ok := err.(*clickhouse.Exception); ok {
			fmt.Printf("Exception [%d] %s \n%s\n", exception.Code, exception.Message, exception.StackTrace)
		}
		return nil, err
	}
	return conn, nil
}

func applyMigrations(config *DatabaseConfig) error {
	dsn := fmt.Sprintf(
		"clickhouse://%s/%s?username=%s&password=%s&x-multi-statement=true",
		config.Addr,
		config.Database,
		config.User,
		config.Password,
	)

	m, err := migrate.New(config.MigrationsDir, dsn)
	if err != nil {
		return fmt.Errorf("error on migrations creation step: %v", err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("error on migrations up step: %v", err)
	}
	return nil
}
