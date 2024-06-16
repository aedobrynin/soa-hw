package main

import (
	"context"
	"errors"
	"flag"
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

	"github.com/aedobrynin/soa-hw/statistics/internal/app"
	"github.com/aedobrynin/soa-hw/statistics/internal/logger"
)

func getConfigPath() string {
	var flagConfigPath string

	flag.StringVar(&flagConfigPath, "c", "./.config.yaml", "path to config file")
	flag.Parse()

	return flagConfigPath
}

func applyMigrations(config *app.DatabaseConfig) error {
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

func connectToDb(ctx context.Context, config *app.DatabaseConfig) (clickhouse.Conn, error) {
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

func main() {
	config, err := app.NewConfig(getConfigPath())
	if err != nil {
		log.Fatal(err)
	}

	err = applyMigrations(&config.Database)
	if err != nil {
		log.Fatal(err)
	}

	clickHouseConn, err := connectToDb(context.Background(), &config.Database)
	if err != nil {
		log.Fatal(err)
	}

	logger, err := logger.GetLogger(config.App.Debug)
	if err != nil {
		log.Fatal(err)
	}

	a := app.New(clickHouseConn, logger, config)

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
