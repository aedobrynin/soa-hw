package repomock

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/aedobrynin/soa-hw/statistics/internal/model"
	"github.com/aedobrynin/soa-hw/statistics/internal/repo"
)

type StatisticsMock struct {
	mock.Mock
}

func (m *StatisticsMock) GetPostStatistics(
	ctx context.Context,
	postID model.PostID,
) (stats *model.PostStatistics, err error) {
	args := m.Called(ctx, postID)
	return args.Get(0).(*model.PostStatistics), args.Error(1)
}

func (m *StatisticsMock) GetTopPosts(
	ctx context.Context,
	request model.GetTopPostsRequest,
) ([]model.CutPostStatistics, error) {
	args := m.Called(ctx, request)
	return args.Get(0).([]model.CutPostStatistics), args.Error(1)
}

func (m *StatisticsMock) GetTopUsersByLikesCount(ctx context.Context, limit uint64) ([]model.UserStatistics, error) {
	args := m.Called(ctx, limit)
	return args.Get(0).([]model.UserStatistics), args.Error(1)
}

func NewPost() repo.Statistics {
	return &StatisticsMock{}
}
