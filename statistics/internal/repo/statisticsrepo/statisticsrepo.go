package statisticsrepo

import (
	"context"
	"fmt"

	"github.com/ClickHouse/clickhouse-go/v2"
	clickhouse_driver "github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/aedobrynin/soa-hw/statistics/internal/model"
	"github.com/aedobrynin/soa-hw/statistics/internal/repo"
	"github.com/google/uuid"
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
	row := r.conn.QueryRow(ctx, "SELECT count(DISTINCT user_id) FROM posts_views WHERE post_id = $1;", postID)
	err := row.Scan(&viewsCnt)
	if err != nil {
		return nil, fmt.Errorf("error on get post views count: %v", err)
	}

	var likesCnt uint64
	row = r.conn.QueryRow(ctx, "SELECT count(DISTINCT user_id) FROM posts_likes WHERE post_id = $1;", postID)
	err = row.Scan(&likesCnt)
	if err != nil {
		return nil, fmt.Errorf("error on get post likes count: %v", err)
	}
	return &model.PostStatistics{LikesCnt: likesCnt, ViewsCnt: viewsCnt}, nil
}

type CutPostStatistics struct {
	PostID   string  `ch:"post_id"`
	LikesCnt *uint64 `ch:"likes_cnt"`
	ViewsCnt *uint64 `ch:"views_cnt"`
}

func (r *statisticsRepo) GetTopPosts(
	ctx context.Context,
	request model.GetTopPostsRequest,
) ([]model.CutPostStatistics, error) {
	var rows clickhouse_driver.Rows
	var err error

	if request.OrderBy == model.OrderByLikesCnt {
		rows, err = r.conn.Query(
			ctx,
			fmt.Sprintf(
				"SELECT post_id, count(DISTINCT user_id) AS likes_cnt FROM posts_likes GROUP BY post_id ORDER BY likes_cnt DESC LIMIT %d;",
				request.Limit,
			),
		)
	} else {
		rows, err = r.conn.Query(
			ctx,
			fmt.Sprintf(
				"SELECT post_id, count(DISTINCT user_id) AS views_cnt FROM posts_views GROUP BY post_id ORDER BY views_cnt DESC LIMIT %d;",
				request.Limit,
			),
		)
	}
	if err != nil {
		return nil, fmt.Errorf("error on get top posts: %v", err)
	}
	defer rows.Close()

	res := make([]model.CutPostStatistics, 0)
	for rows.Next() {
		var cutPostStats CutPostStatistics
		if err := rows.ScanStruct(&cutPostStats); err != nil {
			return nil, fmt.Errorf("error on scanning post stats: %v", err)
		}

		postID, err := uuid.Parse(cutPostStats.PostID)
		if err != nil {
			return nil, fmt.Errorf("error on converting post_id from string to uuid: %v", err)
		}
		res = append(
			res,
			model.CutPostStatistics{PostID: postID, LikesCnt: cutPostStats.LikesCnt, ViewsCnt: cutPostStats.ViewsCnt},
		)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error on get top posts: %v", err)
	}
	return res, nil
}

type UserStatistics struct {
	UserID   string `ch:"post_author_id"`
	LikesCnt uint64 `ch:"likes_cnt"`
}

func (r *statisticsRepo) GetTopUsersByLikesCount(ctx context.Context, limit uint64) ([]model.UserStatistics, error) {
	var rows clickhouse_driver.Rows
	var err error

	rows, err = r.conn.Query(
		ctx,
		fmt.Sprintf(
			"SELECT post_author_id, count(DISTINCT (user_id, post_id)) as likes_cnt FROM posts_likes_indexed_by_post_author GROUP BY post_author_id ORDER BY likes_cnt LIMIT %d",
			limit,
		),
	)

	if err != nil {
		return nil, fmt.Errorf("error on get top users: %v", err)
	}
	defer rows.Close()

	res := make([]model.UserStatistics, 0)
	for rows.Next() {
		var userStats UserStatistics
		if err := rows.ScanStruct(&userStats); err != nil {
			return nil, fmt.Errorf("error on scanning user stats: %v", err)
		}

		userID, err := uuid.Parse(userStats.UserID)
		if err != nil {
			return nil, fmt.Errorf("error on converting user_id from string to uuid: %v", err)
		}
		res = append(
			res,
			model.UserStatistics{UserID: userID, LikesCount: userStats.LikesCnt},
		)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error on get top users: %v", err)
	}
	return res, nil
}

func New(logger *zap.Logger, conn clickhouse.Conn) repo.Statistics {
	return &statisticsRepo{
		logger: logger,
		conn:   conn,
	}
}
