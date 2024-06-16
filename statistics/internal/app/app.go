package app

import (
	"github.com/ClickHouse/clickhouse-go/v2"
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
	if err := a.grpcAdapter.Run(); err != nil {
		defer func() {
			_ = a.logger.Sync()
		}()
		a.logger.Error(err.Error())
	}
	return nil
}

func (a *AppImpl) Shutdown() {
	a.grpcAdapter.Stop()
}

func New(clickHouseConn clickhouse.Conn, logger *zap.Logger, config *Config) App {
	statisticsRepo := statisticsrepo.New(logger, clickHouseConn)
	statisticsService := statisticssvc.New(logger, statisticsRepo)

	return &AppImpl{
		logger:            logger,
		config:            config,
		statisticsService: statisticsService,
		grpcAdapter:       grpcadapter.New(logger, statisticsService, &config.GRPC),
	}
}
