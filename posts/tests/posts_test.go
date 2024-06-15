package tests

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/aedobrynin/soa-hw/posts/internal/app"
	"github.com/aedobrynin/soa-hw/posts/internal/grpcadapter"
	"github.com/aedobrynin/soa-hw/posts/internal/grpcadapter/gen"
	"github.com/aedobrynin/soa-hw/posts/internal/logger"
)

type TestSuite struct {
	suite.Suite
	psqlContainer *PostgreSQLContainer
	app           *app.App
	cfg           *app.Config
	client        gen.PostsClient
}

func (s *TestSuite) SetupSuite() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer ctxCancel()

	psqlContainer, err := NewPostgreSQLContainer(ctx)
	s.Require().NoError(err)

	s.psqlContainer = psqlContainer

	cfg := app.Config{
		App: app.AppConfig{
			Debug: app.DefaultDebug,
		},
		Database: app.DatabaseConfig{
			DSN:           psqlContainer.GetDSN(),
			MigrationsDir: "file://../postgresql/posts/migrations/",
		},
		GRPC: grpcadapter.GRPCConfig{
			Port: app.DefaultGRPCPort,
		},
	}

	m, err := migrate.New(cfg.Database.MigrationsDir, cfg.Database.DSN)
	s.Require().NoError(err)
	err = m.Up()
	s.Require().NoError(err)

	logger, err := logger.GetLogger(cfg.App.Debug)
	s.Require().NoError(err)

	pgxConfig, err := pgxpool.ParseConfig(cfg.Database.DSN)
	s.Require().NoError(err)

	pgxPool, err := pgxpool.ConnectConfig(ctx, pgxConfig)
	s.Require().NoError(err)

	s.app = app.New(pgxPool, logger, &cfg)
	s.cfg = &cfg
}

func (s *TestSuite) TearDownSuite() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	s.Require().NoError(s.psqlContainer.Terminate(ctx))

	s.app.Shutdown()
}

func (s *TestSuite) SetupTest() {
	go func() {
		_ = s.app.Serve()
	}()

	cc, err := grpc.DialContext(context.Background(), fmt.Sprintf(":%d", s.cfg.GRPC.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	s.Require().NoError(err)
	s.client = gen.NewPostsClient(cc)
}

func TestSuite_Run(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (s *TestSuite) TestCreateAndGetPost() {
	authorID := "author"
	content := "content"

	createResp, err := s.client.CreatePost(
		context.Background(),
		&gen.CreatePostRequest{AuthorId: authorID, Content: content},
	)
	s.Require().NoError(err)

	postID := createResp.PostId

	getResp, err := s.client.GetPost(context.Background(), &gen.GetPostRequest{PostId: postID})
	s.Require().NoError(err)
	s.Require().Equal(getResp.Id, postID)
	s.Require().Equal(getResp.AuthorId, authorID)
	s.Require().Equal(getResp.Content, content)
}
