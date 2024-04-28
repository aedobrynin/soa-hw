package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/ClickHouse/clickhouse-go/v2"
	_ "github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/clickhouse"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/aedobrynin/soa-hw/statistics/internal/httpadapter"
)

type app struct {
	config      *Config
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
	// dbConn, err := connectToDb(context.Background(), &config.Database)
	// if err != nil {
	// 	return nil, fmt.Errorf("error on db connection: %v", err)
	// }

	err := applyMigrations(&config.Database)
	if err != nil {
		return nil, fmt.Errorf("error on applyMigrations: %v", err)
	}

	a := &app{
		config:      config,
		httpAdapter: httpadapter.NewAdapter(&config.HTTP),
	}

	return a, nil
}

// func connectToDb(ctx context.Context, config *DatabaseConfig) (driver.Conn, error) {
// 	conn, err := clickhouse.Open(&clickhouse.Options{
// 		Addr: []string{"<CLICKHOUSE_SECURE_NATIVE_HOSTNAME>:9440"},
// 		Auth: clickhouse.Auth{
// 			Database: "default",
// 			Username: "default",
// 			Password: "<DEFAULT_USER_PASSWORD>",
// 		},
// 		ClientInfo: clickhouse.ClientInfo{
// 			Products: []struct {
// 				Name    string
// 				Version string
// 			}{
// 				{Name: "statistics-service-go-client", Version: "0.1"},
// 			},
// 		},

// 		Debugf: func(format string, v ...interface{}) {
// 			log.Printf(format, v)
// 		},
// 		TLS: &tls.Config{
// 			InsecureSkipVerify: true,
// 		},
// 	})

// 	if err != nil {
// 		return nil, err
// 	}

// 	if err := conn.Ping(ctx); err != nil {
// 		if exception, ok := err.(*clickhouse.Exception); ok {
// 			fmt.Printf("Exception [%d] %s \n%s\n", exception.Code, exception.Message, exception.StackTrace)
// 		}
// 		return nil, err
// 	}
// 	return conn, nil
// }

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

	if err := m.Up(); err != nil {
		return fmt.Errorf("error on migrations up step: %v", err)
	}
	return nil
}
