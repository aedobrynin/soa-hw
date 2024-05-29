package statisticsclient

import (
	"context"
	"fmt"

	grpcretry "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/aedobrynin/soa-hw/core/internal/clients"
	"github.com/aedobrynin/soa-hw/core/internal/clients/statisticsclient/gen"
	"github.com/aedobrynin/soa-hw/core/internal/model"
)

type StatisticsClient struct {
	api gen.StatisticsClient
}

var _ clients.StatisticsClient = &StatisticsClient{}

func converToInternal(external *gen.PostStatistics) *model.PostStatistics {
	return &model.PostStatistics{
		PostID:     external.PostId,
		ViewsCount: &external.ViewsCnt,
		LikesCount: &external.LikesCnt,
	}
}

func (c *StatisticsClient) GetPostStatistics(ctx context.Context, postID model.PostID) (*model.PostStatistics, error) {
	stats, err := c.api.GetPostStatistics(ctx, &gen.GetPostStatisticsRequest{PostId: postID})
	if err != nil {
		// TODO: log error
		return nil, fmt.Errorf("error on getting post statistics for id=%s: %v", postID, err)
	}
	return converToInternal(stats), nil
}

func New(
	ctx context.Context,
	config *StatisticsClientConfig,
) (clients.StatisticsClient, error) {
	retryOpts := []grpcretry.CallOption{
		grpcretry.WithCodes(codes.NotFound, codes.Aborted, codes.DeadlineExceeded),
		grpcretry.WithMax(uint(config.RetriesCount)),
		grpcretry.WithPerRetryTimeout(config.Timeout),
	}

	cc, err := grpc.DialContext(ctx, config.Addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(
			grpcretry.UnaryClientInterceptor(retryOpts...),
		))
	if err != nil {
		return nil, fmt.Errorf("error on gRPC connection creation: %w", err)
	}

	grpcClient := gen.NewStatisticsClient(cc)

	return &StatisticsClient{
		api: grpcClient,
	}, nil
}
