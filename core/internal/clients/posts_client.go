package clients

import (
	"context"

	"github.com/google/uuid"

	"github.com/aedobrynin/soa-hw/core/internal/clients/postsclient/gen"
)

type PostsClient interface {
	CreatePost(ctx context.Context, authorId uuid.UUID, content string) (postID uuid.UUID, err error)
	EditPost(ctx context.Context, postId uuid.UUID, editorId uuid.UUID, newContent string) error
	DeletePost(ctx context.Context, postId uuid.UUID, deleterId uuid.UUID) error
	GetPost(ctx context.Context, postId uuid.UUID) (*gen.Post, error)
	ListPosts(
		ctx context.Context,
		pageSize uint32,
		pageToken string,
	) (posts []*gen.Post, nextPageToken string, err error)
}
