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

func (s *statisticsSvc) GetTopPosts(
	ctx context.Context,
	request model.GetTopPostsRequest,
) ([]model.CutPostStatistics, error) {
	defer func() {
		_ = s.logger.Sync()
	}()
	s.logger.Sugar().Debugf("Trying to get top posts. limit=%d, order_by=%d", request.Limit, request.OrderBy)
	if request.Limit > 30 {
		s.logger.Sugar().Infof("Limit=%d is too big", request.Limit)
		return nil, service.ErrLimitTooBig
	}

	top, err := s.repo.GetTopPosts(ctx, request)
	if err != nil {
		wrappedErr := fmt.Errorf("error on getting top posts from repo: %v", err)
		s.logger.Sugar().Error(wrappedErr)
		return nil, wrappedErr
	}
	return top, nil
}

func (s *statisticsSvc) GetTopUsersByLikesCount(ctx context.Context, limit uint64) ([]model.UserStatistics, error) {
	defer func() {
		_ = s.logger.Sync()
	}()
	s.logger.Sugar().Debugf("Trying to get top users. limit=%d", limit)
	if limit > 30 {
		s.logger.Sugar().Infof("Limit=%d is too big", limit)
		return nil, service.ErrLimitTooBig
	}

	top, err := s.repo.GetTopUsersByLikesCount(ctx, limit)
	if err != nil {
		wrappedErr := fmt.Errorf("error on getting top users from repo: %v", err)
		s.logger.Sugar().Error(wrappedErr)
		return nil, wrappedErr
	}
	return top, nil
}

func New(logger *zap.Logger, repo repo.Statistics) service.Statistics {
	return &statisticsSvc{logger: logger, repo: repo}
}
