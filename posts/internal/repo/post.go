package repo

import (
	"context"

	"github.com/aedobrynin/soa-hw/posts/internal/model"
)

type Post interface {
	WithNewTx(ctx context.Context, f func(ctx context.Context) error) error
	AddPost(ctx context.Context, authorID model.UserID, content string) (postID model.PostID, err error)
	GetPost(ctx context.Context, postID model.PostID) (*model.Post, error)
	EditPost(ctx context.Context, postID model.PostID, content string) error
	DeletePost(ctx context.Context, postID model.PostID) error
	// TODO: better
	ListPosts(ctx context.Context, from int, to int) ([]model.Post, error)
}
