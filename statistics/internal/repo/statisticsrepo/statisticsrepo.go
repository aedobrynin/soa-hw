package statisticsrepo

import (
	"context"
	"errors"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/aedobrynin/soa-hw/statistics/internal/model"
	"github.com/aedobrynin/soa-hw/statistics/internal/repo"
	"go.uber.org/zap"
)

type statisticsRepo struct {
	logger *zap.Logger
	conn   clickhouse.Conn
}

func (r *statisticsRepo) GetPostStatistics(
	ctx context.Context,
	postID model.PostID,
) (stats *model.PostStatistics, err error) {
	// TODO
	return nil, errors.New("not implemented")
}

func New(logger *zap.Logger, conn clickhouse.Conn) repo.Statistics {
	return &statisticsRepo{
		logger: logger,
		conn:   conn,
	}
}
