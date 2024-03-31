package postsvc

import (
	"context"

	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/aedobrynin/soa-hw/posts/internal/repo"
	"github.com/aedobrynin/soa-hw/posts/internal/service"
)

type postSvc struct {
	logger *zap.Logger
	repo   repo.Post
}

func (s *postSvc) AddPost(ctx context.Context, authorId uuid.UUID, content string) error {
	// TODO
	return nil
}

func New(logger *zap.Logger, repo repo.Post) service.Post {
	return &postSvc{logger: logger, repo: repo}
}
