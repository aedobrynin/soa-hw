package service

import (
	"context"

	"github.com/aedobrynin/soa-hw/posts/internal/model"
)

type Post interface {
	AddPost(ctx context.Context, authorID model.UserID, content string) (postID model.PostID, err error)
	EditPost(ctx context.Context, postID model.PostID, editorID model.UserID, newContent string) error
	DeletePost(ctx context.Context, postID model.PostID, deleterID model.UserID) error
	GetPost(ctx context.Context, postID model.PostID) (*model.Post, error)
	ListPosts(ctx context.Context, pageSize int, pageToken string) (posts []model.Post, newPageToken string, err error)
}
