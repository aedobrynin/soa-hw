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
	authorId, err := uuid.Parse(request.AuthorId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "author_id should be valid uuid")
	}

	postId, err := s.post.AddPost(ctx, authorId, request.Content)
	if errors.Is(err, service.ErrContentIsEmpty) {
		return nil, status.Errorf(codes.InvalidArgument, "content is empty")
	}
	if err != nil {
		return nil, status.Errorf(codes.Internal, "internal error")
	}
	return &gen.CreatePostResponse{
		PostId: postId.String(),
	}, nil
}

func (s *serverAPI) EditPost(ctx context.Context, request *gen.EditPostRequest) (*empty.Empty, error) {
	postId, err := uuid.Parse(request.PostId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "post_id should be valid uuid")
	}

	editorId, err := uuid.Parse(request.EditorId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "editor_id should be valid uuid")
	}

	err = s.post.EditPost(ctx, postId, editorId, request.NewContent)
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
	postId, err := uuid.Parse(request.PostId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "post_id should be valid uuid")
	}

	deleterId, err := uuid.Parse(request.DeleterId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "deleter_id should be valid uuid")
	}

	err = s.post.DeletePost(ctx, postId, deleterId)
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
	postId, err := uuid.Parse(request.PostId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "post_id should be valid uuid")
	}
	post, err := s.post.GetPost(ctx, postId)
	if errors.Is(err, service.ErrPostNotFound) {
		return nil, status.Errorf(codes.NotFound, "post not found")
	}
	if err != nil {
		return nil, status.Errorf(codes.Internal, "internal error")
	}
	return &gen.Post{
		AuthorId: post.AuthorId.String(),
		Content:  post.Content,
	}, nil
}
