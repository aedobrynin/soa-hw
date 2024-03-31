package postgrpc

import (
	"context"

	"github.com/gofrs/uuid"
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
	if request.AuthorId == "" {
		return nil, status.Error(codes.InvalidArgument, "author_id is required")
	}
	if request.Content == "" {
		return nil, status.Error(codes.InvalidArgument, "content is required")
	}

	authorId, err := uuid.FromString(request.AuthorId)
	if err != nil {
		return &gen.CreatePostResponse{}, status.Error(codes.InvalidArgument, "author_id should be valid uuid")
	}
	postId, err := s.post.AddPost(ctx, authorId, request.Content)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "internal error")
	}
	return &gen.CreatePostResponse{
		PostId: postId.String(),
	}, nil
}
