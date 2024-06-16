package app

import (
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"

	"github.com/aedobrynin/soa-hw/posts/internal/grpcadapter"
	"github.com/aedobrynin/soa-hw/posts/internal/repo/postrepo"
	"github.com/aedobrynin/soa-hw/posts/internal/service"
	"github.com/aedobrynin/soa-hw/posts/internal/service/postsvc"
)

// TODO: interface
type App struct {
	config      *Config
	postService service.Post
	grpcAdapter *grpcadapter.Adapter
	logger      *zap.Logger
}

func (a *App) Serve() error {
	if err := a.grpcAdapter.Run(); err != nil {
		defer func() {
			_ = a.logger.Sync()
		}()
		a.logger.Error(err.Error())
	}
	return nil
}

func (a *App) Shutdown() {
	a.grpcAdapter.Stop()
}

func New(pgxPool *pgxpool.Pool, logger *zap.Logger, config *Config) *App {
	postRepo := postrepo.New(logger, pgxPool)
	postService := postsvc.New(logger, postRepo)

	return &App{
		logger:      logger,
		config:      config,
		postService: postService,
		grpcAdapter: grpcadapter.New(logger, postService, &config.GRPC),
	}
}
