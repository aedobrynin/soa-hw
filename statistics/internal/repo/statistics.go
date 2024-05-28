package repo

import (
	"context"

	"github.com/aedobrynin/soa-hw/statistics/internal/model"
)

type GetTopPostsResponse = []model.CutPostStatistics

type Statistics interface {
	GetPostStatistics(ctx context.Context, postID model.PostID) (stats *model.PostStatistics, err error)

	// If OrderByLikesCnt: LikesCnt != nil, ViewsCnt == nil
	// If OrderByViewsCnt: LikesCnt == nil, ViewsCnt != nil
	GetTopPosts(ctx context.Context, request model.GetTopPostsRequest) ([]model.CutPostStatistics, error)
}
