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

func NewPost() repo.Statistics {
	return &StatisticsMock{}
}
