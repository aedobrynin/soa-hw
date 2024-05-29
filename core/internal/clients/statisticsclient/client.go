package statisticsclient

import (
	"context"
	"errors"
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

func (c *StatisticsClient) GetTopPosts(ctx context.Context, orderBy clients.OrderBy) ([]model.PostStatistics, error) {
	var orderByExternal gen.GetTopPostsRequest_OrderBy
	switch orderBy {
	case clients.OrderByLikesCount:
		orderByExternal = gen.GetTopPostsRequest_LIKES_CNT
	case clients.OrderByViewsCount:
		orderByExternal = gen.GetTopPostsRequest_VIEWS_CNT
	default:
		return nil, errors.New("bad order_by value")
	}
	top, err := c.api.GetTopPosts(ctx, &gen.GetTopPostsRequest{OrderBy: orderByExternal, Limit: 5})
	if err != nil {
		return nil, fmt.Errorf("error on getting top posts from statistics service: %v", err)
	}

	res := make([]model.PostStatistics, 0, len(top.Top))
	for _, postStats := range top.Top {
		res = append(res, model.PostStatistics{
			PostID:     postStats.PostId,
			ViewsCount: postStats.ViewsCnt,
			LikesCount: postStats.LikesCnt,
		})
	}
	return res, nil
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
