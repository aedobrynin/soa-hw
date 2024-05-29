package clients

import (
	"context"

	"github.com/aedobrynin/soa-hw/core/internal/model"
)

type StatisticsClient interface {
	GetPostStatistics(ctx context.Context, postID model.PostID) (*model.PostStatistics, error)
}
