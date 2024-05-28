package statisticsgrpc

import (
	"context"
	"errors"

	"google.golang.org/grpc"

	"github.com/aedobrynin/soa-hw/statistics/internal/grpcadapter/gen"
	"github.com/aedobrynin/soa-hw/statistics/internal/service"
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
	// TODO
	return nil, errors.New("not implemented")
}
