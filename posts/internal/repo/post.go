package repo

import (
	"context"

	"github.com/google/uuid"

	"github.com/aedobrynin/soa-hw/posts/internal/model"
)

type Post interface {
	WithNewTx(ctx context.Context, f func(ctx context.Context) error) error
	AddPost(ctx context.Context, authorId uuid.UUID, content string) (postId uuid.UUID, err error)
	GetPost(ctx context.Context, postId uuid.UUID) (*model.Post, error)
	EditPost(ctx context.Context, postId uuid.UUID, content string) error
	DeletePost(ctx context.Context, postId uuid.UUID) error
	// TODO: better
	ListPosts(ctx context.Context, from int, to int) ([]model.Post, error)
}
