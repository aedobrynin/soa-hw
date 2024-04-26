package postgrpc

import (
	"context"
	"errors"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/aedobrynin/soa-hw/posts/internal/grpcadapter/gen"
	"github.com/aedobrynin/soa-hw/posts/internal/service"
)

type serverAPI struct {
	gen.UnimplementedPostsServer
	post service.Post
}

func Register(gRPCServer *grpc.Server, post service.Post) {
	gen.RegisterPostsServer(gRPCServer, &serverAPI{post: post})
}

func (s *serverAPI) CreatePost(
	ctx context.Context,
	request *gen.CreatePostRequest,
) (*gen.CreatePostResponse, error) {
	postID, err := s.post.AddPost(ctx, request.AuthorId, request.Content)
	if errors.Is(err, service.ErrContentIsEmpty) {
		return nil, status.Errorf(codes.InvalidArgument, "content is empty")
	}
	if err != nil {
		return nil, status.Errorf(codes.Internal, "internal error")
	}
	return &gen.CreatePostResponse{
		PostId: postID.String(),
	}, nil
}

func (s *serverAPI) EditPost(ctx context.Context, request *gen.EditPostRequest) (*empty.Empty, error) {
	postID, err := uuid.Parse(request.PostId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "post_id should be valid uuid")
	}

	err = s.post.EditPost(ctx, postID, request.EditorId, request.NewContent)
	if errors.Is(err, service.ErrContentIsEmpty) {
		return nil, status.Errorf(codes.InvalidArgument, "new_content is empty")
	}
	if errors.Is(err, service.ErrPostNotFound) {
		return nil, status.Errorf(codes.NotFound, "post not found")
	}
	if errors.Is(err, service.ErrInsufficientPermissions) {
		return nil, status.Errorf(codes.PermissionDenied, "only post creator can edit post")
	}
	if err != nil {
		return nil, status.Errorf(codes.Internal, "internal error")
	}
	return &empty.Empty{}, nil
}

func (s *serverAPI) DeletePost(ctx context.Context, request *gen.DeletePostRequest) (*empty.Empty, error) {
	postID, err := uuid.Parse(request.PostId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "post_id should be valid uuid")
	}

	err = s.post.DeletePost(ctx, postID, request.DeleterId)
	if errors.Is(err, service.ErrPostNotFound) {
		return nil, status.Errorf(codes.NotFound, "post not found")
	}
	if errors.Is(err, service.ErrInsufficientPermissions) {
		return nil, status.Errorf(codes.PermissionDenied, "only post creator can edit post")
	}
	if err != nil {
		return nil, status.Errorf(codes.Internal, "internal error")
	}
	return &empty.Empty{}, nil
}

func (s *serverAPI) GetPost(ctx context.Context, request *gen.GetPostRequest) (*gen.Post, error) {
	postID, err := uuid.Parse(request.PostId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "post_id should be valid uuid")
	}
	post, err := s.post.GetPost(ctx, postID)
	if errors.Is(err, service.ErrPostNotFound) {
		return nil, status.Errorf(codes.NotFound, "post not found")
	}
	if err != nil {
		return nil, status.Errorf(codes.Internal, "internal error")
	}
	return &gen.Post{
		Id:       post.ID.String(),
		AuthorId: post.AuthorID,
		Content:  post.Content,
	}, nil
}

func (s *serverAPI) ListPosts(ctx context.Context, request *gen.ListPostsRequest) (*gen.ListPostsResponse, error) {
	posts, nextPageToken, err := s.post.ListPosts(ctx, int(request.PageSize), request.PageToken)
	if errors.Is(err, service.ErrBadPageToken) {
		return nil, status.Errorf(codes.InvalidArgument, "bad page token")
	}
	if err != nil {
		return nil, status.Errorf(codes.Internal, "internal error")
	}

	respPosts := make([]*gen.Post, 0, len(posts))
	for _, post := range posts {
		respPosts = append(
			respPosts,
			&gen.Post{Id: post.ID.String(), AuthorId: post.AuthorID, Content: post.Content},
		)
	}
	return &gen.ListPostsResponse{
		Posts:         respPosts,
		NextPageToken: nextPageToken,
	}, nil
}
