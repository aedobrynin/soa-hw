package statisticsgrpc

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/aedobrynin/soa-hw/statistics/internal/grpcadapter/gen"
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
