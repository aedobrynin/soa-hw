package postsclient

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	grpcretry "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
	"github.com/juju/zaputil/zapctx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	"github.com/aedobrynin/soa-hw/core/internal/clients"
	"github.com/aedobrynin/soa-hw/core/internal/clients/postsclient/gen"
)

type PostsClient struct {
	api gen.PostsClient
}

func (c *PostsClient) CreatePost(
	ctx context.Context,
	authorId uuid.UUID,
	content string,
) (uuid.UUID, error) {
	resp, err := c.api.CreatePost(ctx, &gen.CreatePostRequest{AuthorId: authorId.String(), Content: content})
	if err != nil {
		// TODO: this is bad
		if status.Code(err) == codes.InvalidArgument {
			return uuid.Nil, clients.ErrContentIsEmpty
		}
		return uuid.Nil, err
	}

	postId, err := uuid.Parse(resp.PostId)
	if err != nil {
		logger := zapctx.Logger(ctx)
		defer logger.Sync()
		logger.Sugar().Errorf("couldn't parse postId %s from posts service as uuid", resp.PostId)
		return uuid.Nil, err
	}
	return postId, nil
}

func (c *PostsClient) EditPost(
	ctx context.Context,
	postId uuid.UUID,
	editorId uuid.UUID,
	newContent string,
) error {
	_, err := c.api.EditPost(
		ctx,
		&gen.EditPostRequest{PostId: postId.String(), EditorId: editorId.String(), NewContent: newContent},
	)
	if err != nil {
		// TODO: this is bad
		if status.Code(err) == codes.InvalidArgument {
			return clients.ErrContentIsEmpty
		}
		if status.Code(err) == codes.NotFound {
			return clients.ErrPostNotFound
		}
		if status.Code(err) == codes.PermissionDenied {
			return clients.ErrInsufficientPermissions
		}
		return err
	}
	return nil
}

func (c *PostsClient) DeletePost(
	ctx context.Context,
	postId uuid.UUID,
	deleterId uuid.UUID,
) error {
	_, err := c.api.DeletePost(
		ctx,
		&gen.DeletePostRequest{PostId: postId.String(), DeleterId: deleterId.String()},
	)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return clients.ErrPostNotFound
		}
		if status.Code(err) == codes.PermissionDenied {
			return clients.ErrInsufficientPermissions
		}
		return err
	}
	return nil
}

func (c *PostsClient) GetPost(
	ctx context.Context,
	postId uuid.UUID,
) (*gen.Post, error) {
	post, err := c.api.GetPost(ctx, &gen.GetPostRequest{PostId: postId.String()})
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, clients.ErrPostNotFound
		}
		return nil, err
	}
	return post, nil
}

func New(
	ctx context.Context,
	config *PostsClientConfig,
) (*PostsClient, error) {
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

	grpcClient := gen.NewPostsClient(cc)

	return &PostsClient{
		api: grpcClient,
	}, nil
}
