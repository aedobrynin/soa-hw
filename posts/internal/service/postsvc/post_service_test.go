package postsvc_test

import (
	"context"
	"testing"

	"github.com/gofrs/uuid"
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

	authorId := uuid.Must(uuid.NewV4())
	content := "content"
	expectedPostId := uuid.Must(uuid.NewV4())

	postRepo := repomock.NewPost()
	postRepo.On(
		"AddPost",
		mock.MatchedBy(func(ctx context.Context) bool { return true }),
		authorId,
		content,
	).Return(expectedPostId, nil)
	svc := postsvc.New(logger, postRepo)

	ctx := context.Background()

	returnedPostId, err := svc.AddPost(ctx, authorId, "content")
	require.Equal(t, returnedPostId, expectedPostId, "returned different post id")
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

	returnedPostId, err := svc.AddPost(ctx, uuid.Must(uuid.NewV4()), "")
	require.Equal(t, returnedPostId, uuid.Nil)
	require.Equal(t, err, service.ErrContentIsEmpty)
}
