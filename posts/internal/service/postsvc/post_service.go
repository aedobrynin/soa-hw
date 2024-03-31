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

func (s *postSvc) AddPost(ctx context.Context, authorId uuid.UUID, content string) (uuid.UUID, error) {
	defer s.logger.Sync()
	s.logger.Sugar().Infof("Trying to add post with author_id=%s, content=%s", authorId.String(), content)

	if len(content) == 0 {
		return uuid.Nil, service.ErrContentIsEmpty
	}

	return s.repo.AddPost(ctx, authorId, content)
}

func New(logger *zap.Logger, repo repo.Post) service.Post {
	return &postSvc{logger: logger, repo: repo}
}
