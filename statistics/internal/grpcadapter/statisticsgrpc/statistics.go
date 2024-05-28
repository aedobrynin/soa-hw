package statisticsgrpc

import (
	"context"
	"errors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/aedobrynin/soa-hw/statistics/internal/grpcadapter/gen"
	"github.com/aedobrynin/soa-hw/statistics/internal/model"
	"github.com/aedobrynin/soa-hw/statistics/internal/service"
	"github.com/google/uuid"
)

type serverAPI struct {
	gen.UnimplementedStatisticsServer
	statistics service.Statistics
}

func Register(gRPCServer *grpc.Server, statistics service.Statistics) {
	gen.RegisterStatisticsServer(gRPCServer, &serverAPI{statistics: statistics})
}

func (s *serverAPI) GetPostStatistics(
	ctx context.Context,
	request *gen.GetPostStatisticsRequest,
) (*gen.PostStatistics, error) {
	postID, err := uuid.Parse(request.PostId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "post_id should be valid uuid")
	}

	stats, err := s.statistics.GetPostStatistics(ctx, postID)
	if err != nil {

		return nil, status.Errorf(codes.Internal, "internal error")
	}
	return &gen.PostStatistics{
		PostId:   request.PostId,
		LikesCnt: stats.LikesCnt,
		ViewsCnt: stats.ViewsCnt,
	}, nil
}

func (s *serverAPI) GetTopPosts(
	ctx context.Context,
	request *gen.GetTopPostsRequest,
) (*gen.GetTopPostsResponse, error) {
	var orderBy model.OrderBy
	switch request.OrderBy {
	case gen.GetTopPostsRequest_LIKES_CNT:
		orderBy = model.OrderByLikesCnt
	case gen.GetTopPostsRequest_VIEWS_CNT:
		orderBy = model.OrderByViewsCnt
	default:
		return nil, status.Error(codes.InvalidArgument, "bad order_by value")
	}

	top, err := s.statistics.GetTopPosts(
		ctx,
		model.GetTopPostsRequest{Limit: request.Limit, OrderBy: orderBy},
	)
	if errors.Is(err, service.ErrLimitTooBig) {
		return nil, status.Error(codes.InvalidArgument, "limit is too big")
	}
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	res := make([]*gen.CutPostStatistics, 0, len(top))
	for _, postStats := range top {
		res = append(
			res,
			&gen.CutPostStatistics{
				PostId:   postStats.PostID.String(),
				LikesCnt: postStats.LikesCnt,
				ViewsCnt: postStats.ViewsCnt,
			},
		)
	}
	return &gen.GetTopPostsResponse{Top: res}, nil
}
