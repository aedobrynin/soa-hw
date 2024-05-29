package clients

import (
	"context"

	"github.com/aedobrynin/soa-hw/core/internal/model"
)

// TODO: better
type OrderBy = uint8

const (
	OrderByLikesCount OrderBy = 0
	OrderByViewsCount OrderBy = 1
)

type StatisticsClient interface {
	GetPostStatistics(ctx context.Context, postID model.PostID) (*model.PostStatistics, error)
	GetTopPosts(ctx context.Context, orderBy OrderBy) ([]model.PostStatistics, error)
}
