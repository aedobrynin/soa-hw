package tests

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	_ "github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/clickhouse"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/aedobrynin/soa-hw/statistics/internal/app"
	"github.com/aedobrynin/soa-hw/statistics/internal/grpcadapter"
	"github.com/aedobrynin/soa-hw/statistics/internal/grpcadapter/gen"
	"github.com/aedobrynin/soa-hw/statistics/internal/logger"
)

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

type TestSuite struct {
	suite.Suite
	clickHouseContainer *ClickHouseContainer
	app                 app.App
	cfg                 *app.Config
	client              gen.StatisticsClient
}

func (s *TestSuite) SetupSuite() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer ctxCancel()

	clickHouseContainer, err := NewClickHouseContainer(ctx)
	s.Require().NoError(err)

	s.clickHouseContainer = clickHouseContainer

	cfg := app.Config{
		App: app.AppConfig{
			Debug: app.DefaultDebug,
		},
		Database: app.DatabaseConfig{
			Addr:          clickHouseContainer.GetAddr(),
			Database:      clickHouseContainer.Config.Database,
			User:          clickHouseContainer.Config.User,
			Password:      clickHouseContainer.Config.Password,
			MigrationsDir: "file://../clickhouse/statistics/migrations/",
		},
		GRPC: grpcadapter.GRPCConfig{
			Port: app.DefaultGRPCPort,
		},
	}

	migrationsDSN := fmt.Sprintf(
		"clickhouse://%s/%s?username=%s&password=%s&x-multi-statement=true",
		cfg.Database.Addr,
		cfg.Database.Database,
		cfg.Database.User,
		cfg.Database.Password,
	)

	m, err := migrate.New(cfg.Database.MigrationsDir, migrationsDSN)
	s.Require().NoError(err)
	err = m.Up()
	s.Require().NoError(err)

	logger, err := logger.GetLogger(cfg.App.Debug)
	s.Require().NoError(err)

	clickHouseConn, err := connectToDb(context.Background(), &cfg.Database)
	s.Require().NoError(err)

	s.app = app.New(clickHouseConn, logger, &cfg)
	s.cfg = &cfg
}

func (s *TestSuite) TearDownSuite() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	s.Require().NoError(s.clickHouseContainer.Terminate(ctx))

	s.app.Shutdown()
}

func (s *TestSuite) SetupTest() {
	go func() {
		_ = s.app.Serve()
	}()

	cc, err := grpc.DialContext(context.Background(), fmt.Sprintf(":%d", s.cfg.GRPC.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	s.Require().NoError(err)
	s.client = gen.NewStatisticsClient(cc)
}

func TestSuite_Run(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (s *TestSuite) TestGetPostStatistics() {
	clickHouseConn, err := connectToDb(context.Background(), &s.cfg.Database)
	s.Require().NoError(err)

	userID := "00112233-4455-6677-8899-aabbccddeeff"
	postID := "123e4567-e89b-12d3-a456-426655440000"

	err = clickHouseConn.Exec(context.Background(), fmt.Sprintf(`INSERT INTO posts_likes (user_id, post_id)
				  												 VALUES ('%s', '%s')`, userID, postID))
	s.Require().NoError(err)
	err = clickHouseConn.Exec(context.Background(), fmt.Sprintf(`INSERT INTO posts_views (user_id, post_id)
				  												 VALUES ('%s', '%s')`, userID, postID))
	s.Require().NoError(err)

	resp, err := s.client.GetPostStatistics(context.Background(), &gen.GetPostStatisticsRequest{PostId: postID})
	s.Require().NoError(err)

	s.Require().Equal(resp.PostId, postID)
	s.Require().Equal(resp.LikesCnt, uint64(1))
	s.Require().Equal(resp.ViewsCnt, uint64(1))
	_ = clickHouseConn.Close()
}
