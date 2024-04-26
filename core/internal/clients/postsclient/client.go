package postsclient

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	grpcretry "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	"github.com/aedobrynin/soa-hw/core/internal/clients"
	"github.com/aedobrynin/soa-hw/core/internal/clients/postsclient/gen"
	"github.com/aedobrynin/soa-hw/core/internal/model"
)

type PostsClient struct {
	api gen.PostsClient
}

var _ clients.PostsClient = &PostsClient{}

func convertToInternal(post *gen.Post) (*model.Post, error) {
	if post == nil {
		return nil, nil
	}
	authorID, err := uuid.Parse(post.AuthorId)
	if err != nil {
		return nil, err
	}

	return &model.Post{
		ID:       post.Id,
		AuthorID: authorID,
		Content:  post.Content,
	}, nil
}

func (c *PostsClient) CreatePost(
	ctx context.Context,
	authorID model.UserID,
	content string,
) (model.PostID, error) {
	resp, err := c.api.CreatePost(ctx, &gen.CreatePostRequest{AuthorId: authorID.String(), Content: content})
	if err != nil {
		// TODO: is this bad?
		if status.Code(err) == codes.InvalidArgument {
			return "", clients.ErrContentIsEmpty
		}
		return "", err
	}
	return resp.PostId, nil
}

func (c *PostsClient) EditPost(
	ctx context.Context,
	postID model.PostID,
	editorID model.UserID,
	newContent string,
) error {
	_, err := c.api.EditPost(
		ctx,
		&gen.EditPostRequest{PostId: postID, EditorId: editorID.String(), NewContent: newContent},
	)
	if err != nil {
		// TODO: is this bad?
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
	postID model.PostID,
	deleterID model.UserID,
) error {
	_, err := c.api.DeletePost(
		ctx,
		&gen.DeletePostRequest{PostId: postID, DeleterId: deleterID.String()},
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
	postID model.PostID,
) (*model.Post, error) {
	post, err := c.api.GetPost(ctx, &gen.GetPostRequest{PostId: postID})
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, clients.ErrPostNotFound
		}
		return nil, err
	}
	return convertToInternal(post)
}

func (c *PostsClient) ListPosts(
	ctx context.Context,
	pageSize uint32,
	pageToken string,
) (posts []*model.Post, nextPageToken string, err error) {
	resp, err := c.api.ListPosts(ctx, &gen.ListPostsRequest{PageSize: pageSize, PageToken: pageToken})
	if err != nil {
		if status.Code(err) == codes.InvalidArgument {
			return nil, "", clients.ErrBadPageToken
		}
		return nil, "", err
	}

	posts = make([]*model.Post, 0, len(resp.Posts))
	for _, post := range resp.Posts {
		convertedPost, err := convertToInternal(post)
		if err != nil {
			return nil, "", err
		}
		posts = append(posts, convertedPost)
	}
	return posts, resp.NextPageToken, nil
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
