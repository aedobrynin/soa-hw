package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/aedobrynin/soa-hw/posts/internal/model"
)

type Post interface {
	AddPost(ctx context.Context, authorId uuid.UUID, content string) (postId uuid.UUID, err error)
	EditPost(ctx context.Context, postId uuid.UUID, editorId uuid.UUID, newContent string) error
	DeletePost(ctx context.Context, postId uuid.UUID, deleterId uuid.UUID) error
	GetPost(ctx context.Context, postId uuid.UUID) (*model.Post, error)
	ListPosts(ctx context.Context, pageSize int, pageToken string) (posts []model.Post, newPageToken string, err error)
}
