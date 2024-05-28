package service

import (
	"context"

	"github.com/aedobrynin/soa-hw/statistics/internal/model"
)

type Statistics interface {
	GetPostStatistics(ctx context.Context, postId model.PostID) (stats *model.PostStatistics, err error)
}
