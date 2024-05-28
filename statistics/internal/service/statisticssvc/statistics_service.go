package statisticssvc

import (
	"context"
	"errors"

	"go.uber.org/zap"

	"github.com/aedobrynin/soa-hw/statistics/internal/model"
	"github.com/aedobrynin/soa-hw/statistics/internal/repo"
	"github.com/aedobrynin/soa-hw/statistics/internal/service"
)

type statisticsSvc struct {
	logger *zap.Logger
	repo   repo.Statistics
}

var _ service.Statistics = &statisticsSvc{}

func (s *statisticsSvc) GetPostStatistics(
	ctx context.Context,
	postId model.PostID,
) (stats *model.PostStatistics, err error) {
	defer func() {
		_ = s.logger.Sync()
	}()
	// TODO
	return nil, errors.New("not implemented")
}

func New(logger *zap.Logger, repo repo.Statistics) service.Statistics {
	return &statisticsSvc{logger: logger, repo: repo}
}
