package clients

import (
	"context"

	"github.com/aedobrynin/soa-hw/core/internal/model"
)

type PostsClient interface {
	CreatePost(ctx context.Context, authorID model.UserID, content string) (postID string, err error)
	EditPost(ctx context.Context, postID string, editorID model.UserID, newContent string) error
	DeletePost(ctx context.Context, postID string, deleterID model.UserID) error
	GetPost(ctx context.Context, postID string) (*model.Post, error)
	ListPosts(
		ctx context.Context,
		pageSize uint32,
		pageToken string,
	) (posts []*model.Post, nextPageToken string, err error)
}
