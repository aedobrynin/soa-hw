package postsvc_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/aedobrynin/soa-hw/posts/internal/logger"
	"github.com/aedobrynin/soa-hw/posts/internal/repo/repomock"
	"github.com/aedobrynin/soa-hw/posts/internal/service"
	"github.com/aedobrynin/soa-hw/posts/internal/service/postsvc"
)

func TestAddPostHappyPath(t *testing.T) {
	logger, err := logger.GetLogger(true)
	if err != nil {
		t.Error(err)
	}

	authorID := uuid.New().String()
	content := "content"
	expectedPostID := uuid.New()

	postRepo := repomock.NewPost()
	postRepo.On(
		"AddPost",
		mock.MatchedBy(func(ctx context.Context) bool { return true }),
		authorID,
		content,
	).Return(expectedPostID, nil)
	svc := postsvc.New(logger, postRepo)

	ctx := context.Background()

	returnedPostID, err := svc.AddPost(ctx, authorID, "content")
	require.Equal(t, returnedPostID, expectedPostID, "returned different post id")
	require.Nil(t, err, "error should be nil in happy path")
}

func TestAddPostEmptyContent(t *testing.T) {
	logger, err := logger.GetLogger(true)
	if err != nil {
		t.Error(err)
	}

	postRepo := repomock.NewPost()
	svc := postsvc.New(logger, postRepo)

	ctx := context.Background()

	returnedPostID, err := svc.AddPost(ctx, uuid.New().String(), "")
	require.Equal(t, returnedPostID, uuid.Nil)
	require.Equal(t, err, service.ErrContentIsEmpty)
}

func TestEditPostPostNotFound(t *testing.T) {
	// TODO
}

func TestEditPostEmptyContent(t *testing.T) {
	// TODO
}

func TestEditPostInsufficientPermissions(t *testing.T) {
	//TODO
}

func TestEditPostHappyPath(t *testing.T) {
	// TODO
}

func TestDeletePostPostNotFound(t *testing.T) {
	// TODO
}

func TestDeletePostinsufficientPermissions(t *testing.T) {
	// TODO
}

func TestDeletePostHappyPath(t *testing.T) {
	// TODO
}

func TestGetPostPostNotFound(t *testing.T) {
	// TODO
}

func TestGetPostHappyPath(t *testing.T) {
	// TODO
}

// TODO: TestListPosts
