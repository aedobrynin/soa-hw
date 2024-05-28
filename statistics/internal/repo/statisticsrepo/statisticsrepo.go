package statisticsrepo

import (
	"context"
	"fmt"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/aedobrynin/soa-hw/statistics/internal/model"
	"github.com/aedobrynin/soa-hw/statistics/internal/repo"
	"go.uber.org/zap"
)

type statisticsRepo struct {
	logger *zap.Logger
	conn   clickhouse.Conn
}

func (r *statisticsRepo) GetPostStatistics(
	ctx context.Context,
	postID model.PostID,
) (*model.PostStatistics, error) {
	var viewsCnt uint64
	row := r.conn.QueryRow(ctx, "SELECT COUNT(DISTINCT user_id) FROM posts_views WHERE post_id = $1;", postID)
	err := row.Scan(&viewsCnt)
	if err != nil {
		return nil, fmt.Errorf("error on get post views count: %v", err)
	}

	var likesCnt uint64
	row = r.conn.QueryRow(ctx, "SELECT COUNT(DISTINCT user_id) FROM posts_likes WHERE post_id = $1;", postID)
	err = row.Scan(&likesCnt)
	if err != nil {
		return nil, fmt.Errorf("error on get post likes count: %v", err)
	}
	return &model.PostStatistics{LikesCnt: likesCnt, ViewsCnt: viewsCnt}, nil
}

func New(logger *zap.Logger, conn clickhouse.Conn) repo.Statistics {
	return &statisticsRepo{
		logger: logger,
		conn:   conn,
	}
}
