package clients

import (
	"context"

	"github.com/aedobrynin/soa-hw/core/internal/model"
)

type PostsClient interface {
	CreatePost(ctx context.Context, authorId model.UserId, content string) (postID string, err error)
	EditPost(ctx context.Context, postId string, editorId model.UserId, newContent string) error
	DeletePost(ctx context.Context, postId string, deleterId model.UserId) error
	GetPost(ctx context.Context, postId string) (*model.Post, error)
	ListPosts(
		ctx context.Context,
		pageSize uint32,
		pageToken string,
	) (posts []*model.Post, nextPageToken string, err error)
}
