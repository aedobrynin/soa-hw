package postsvc

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/aedobrynin/soa-hw/posts/internal/model"
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

	postId, err := s.repo.AddPost(ctx, authorId, content)
	if err != nil {
		s.logger.Sugar().Infof("Couldn't create post: %s", err)
	} else {
		s.logger.Sugar().Infof("Created post with id=%s", postId)
	}
	return postId, err
}

func (s *postSvc) EditPost(ctx context.Context, postId uuid.UUID, editorId uuid.UUID, newContent string) error {
	if len(newContent) == 0 {
		return service.ErrContentIsEmpty
	}

	post, err := s.repo.GetPost(ctx, postId)
	if errors.Is(err, repo.ErrPostNotFound) {
		return service.ErrPostNotFound
	}
	if err != nil {
		return err
	}

	if post.AuthorId != editorId {
		return service.ErrInsufficientPermissions
	}

	err = s.repo.EditPost(ctx, postId, newContent)
	if errors.Is(err, repo.ErrPostNotFound) {
		return service.ErrPostNotFound
	}
	return err
}

func (s *postSvc) DeletePost(ctx context.Context, postId uuid.UUID, deleterId uuid.UUID) error {
	post, err := s.repo.GetPost(ctx, postId)
	if errors.Is(err, repo.ErrPostNotFound) {
		return service.ErrPostNotFound
	}
	if err != nil {
		return err
	}

	if post.AuthorId != deleterId {
		return service.ErrInsufficientPermissions
	}

	err = s.repo.DeletePost(ctx, postId)
	if errors.Is(err, repo.ErrPostNotFound) {
		return service.ErrPostNotFound
	}
	return err
}
func (s *postSvc) GetPost(ctx context.Context, postId uuid.UUID) (*model.Post, error) {
	post, err := s.repo.GetPost(ctx, postId)
	if errors.Is(err, repo.ErrPostNotFound) {
		return nil, service.ErrPostNotFound
	}
	if err != nil {
		return nil, err
	}

	return post, nil
}

func New(logger *zap.Logger, repo repo.Post) service.Post {
	return &postSvc{logger: logger, repo: repo}
}
