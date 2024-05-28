package statisticssvc

import (
	"context"
	"fmt"

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
	postID model.PostID,
) (*model.PostStatistics, error) {
	defer func() {
		_ = s.logger.Sync()
	}()
	s.logger.Sugar().Debugf("Trying to get post statistics for id=%s", postID)
	stats, err := s.repo.GetPostStatistics(ctx, postID)
	if err != nil {
		wrappedErr := fmt.Errorf("error on getting post stats from repo: %v", err)
		s.logger.Sugar().Error(wrappedErr)
		return nil, wrappedErr
	}
	return stats, nil
}

func New(logger *zap.Logger, repo repo.Statistics) service.Statistics {
	return &statisticsSvc{logger: logger, repo: repo}
}
