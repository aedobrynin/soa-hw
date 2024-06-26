package postsvc

import (
	"context"
	"errors"
	"strconv"

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

func (s *postSvc) AddPost(ctx context.Context, authorID model.UserID, content string) (model.PostID, error) {
	defer func() {
		_ = s.logger.Sync()
	}()
	s.logger.Sugar().Infof("Trying to add post with author_id=%s, content=%s", authorID, content)

	if len(content) == 0 {
		return uuid.Nil, service.ErrContentIsEmpty
	}

	postID, err := s.repo.AddPost(ctx, authorID, content)
	if err != nil {
		s.logger.Sugar().Infof("Couldn't create post: %s", err)
	} else {
		s.logger.Sugar().Infof("Created post with id=%s", postID)
	}
	return postID, err
}

func (s *postSvc) EditPost(ctx context.Context, postID model.PostID, editorID model.UserID, newContent string) error {
	if len(newContent) == 0 {
		return service.ErrContentIsEmpty
	}

	post, err := s.repo.GetPost(ctx, postID)
	if errors.Is(err, repo.ErrPostNotFound) {
		return service.ErrPostNotFound
	}
	if err != nil {
		return err
	}

	if post.AuthorID != editorID {
		return service.ErrInsufficientPermissions
	}

	err = s.repo.EditPost(ctx, postID, newContent)
	if errors.Is(err, repo.ErrPostNotFound) {
		return service.ErrPostNotFound
	}
	return err
}

func (s *postSvc) DeletePost(ctx context.Context, postID model.PostID, deleterID model.UserID) error {
	post, err := s.repo.GetPost(ctx, postID)
	if errors.Is(err, repo.ErrPostNotFound) {
		return service.ErrPostNotFound
	}
	if err != nil {
		return err
	}

	if post.AuthorID != deleterID {
		return service.ErrInsufficientPermissions
	}

	err = s.repo.DeletePost(ctx, postID)
	if errors.Is(err, repo.ErrPostNotFound) {
		return service.ErrPostNotFound
	}
	return err
}
func (s *postSvc) GetPost(ctx context.Context, postID model.PostID) (*model.Post, error) {
	post, err := s.repo.GetPost(ctx, postID)
	if errors.Is(err, repo.ErrPostNotFound) {
		return nil, service.ErrPostNotFound
	}
	if err != nil {
		return nil, err
	}

	return post, nil
}

// TODO: better
func pageTokenToPage(pageToken string) (int, error) {
	if pageToken == "" {
		return 0, nil
	}
	page, err := strconv.Atoi(pageToken)
	if err != nil {
		return 0, service.ErrBadPageToken
	}
	return page, nil
}

func (s *postSvc) ListPosts(ctx context.Context, pageSize int, pageToken string) ([]model.Post, string, error) {
	page, err := pageTokenToPage(pageToken)
	if err != nil {
		return nil, "", err
	}

	if pageSize == 0 {
		pageSize = 5
	}
	// TODO: better
	if pageSize > 100 {
		pageSize = 100
	}

	s.logger.Debug("repo.ListPosts")
	posts, err := s.repo.ListPosts(ctx, page, page+pageSize)
	if err != nil {
		return nil, "", err
	}
	return posts, strconv.Itoa(page + len(posts)), nil
}

func New(logger *zap.Logger, repo repo.Post) service.Post {
	return &postSvc{logger: logger, repo: repo}
}
