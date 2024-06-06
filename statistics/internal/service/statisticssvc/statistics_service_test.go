package statisticssvc_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/aedobrynin/soa-hw/statistics/internal/logger"
	"github.com/aedobrynin/soa-hw/statistics/internal/model"
	"github.com/aedobrynin/soa-hw/statistics/internal/repo/repomock"
	"github.com/aedobrynin/soa-hw/statistics/internal/service/statisticssvc"
)

func TestGetPostStatisticsHappyPath(t *testing.T) {
	logger, err := logger.GetLogger(true)
	if err != nil {
		t.Error(err)
	}

	postID := uuid.New()
	stats := &model.PostStatistics{LikesCnt: 100, ViewsCnt: 500}

	statisticsRepo := repomock.NewStatistics()
	statisticsRepo.On(
		"GetPostStatistics",
		mock.MatchedBy(func(ctx context.Context) bool { return true }),
		postID,
	).Return(stats, nil)
	svc := statisticssvc.New(logger, &statisticsRepo)

	ctx := context.Background()

	returnedStats, err := svc.GetPostStatistics(ctx, postID)
	require.Equal(t, returnedStats, stats, "returned different stats")
	require.Nil(t, err, "error should be nil in happy path")
}

func TestGetPostStatisticsRepoError(t *testing.T) {
	logger, err := logger.GetLogger(true)
	if err != nil {
		t.Error(err)
	}

	postID := uuid.New()

	var nilStats *model.PostStatistics = nil

	statisticsRepo := repomock.NewStatistics()
	statisticsRepo.On(
		"GetPostStatistics",
		mock.MatchedBy(func(ctx context.Context) bool { return true }),
		postID,
	).Return(nilStats, errors.New("some repo error"))
	svc := statisticssvc.New(logger, &statisticsRepo)

	ctx := context.Background()

	returnedStats, err := svc.GetPostStatistics(ctx, postID)
	require.Nil(t, returnedStats, nil, "should return nil stats on repo error")
	require.Error(t, err, "should return not nil error on repo error")
}
