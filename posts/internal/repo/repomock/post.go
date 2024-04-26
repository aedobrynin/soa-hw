package repomock

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/aedobrynin/soa-hw/posts/internal/model"
)

type PostMock struct {
	mock.Mock
}

func (m *PostMock) WithNewTx(ctx context.Context, f func(ctx context.Context) error) error {
	args := m.Called(ctx, f)
	return args.Error(0)
}

func (m *PostMock) AddPost(ctx context.Context, authorID model.UserID, content string) (model.PostID, error) {
	args := m.Called(ctx, authorID, content)
	return args.Get(0).(model.PostID), args.Error(1)
}

func (m *PostMock) GetPost(ctx context.Context, postID model.PostID) (*model.Post, error) {
	args := m.Called(ctx, postID)
	return args.Get(0).(*model.Post), args.Error(1)
}

func (m *PostMock) EditPost(ctx context.Context, postID model.PostID, content string) error {
	args := m.Called(ctx, postID, content)
	return args.Error(0)
}

func (m *PostMock) DeletePost(ctx context.Context, postID model.PostID) error {
	args := m.Called(ctx, postID)
	return args.Error(0)
}

func (m *PostMock) ListPosts(ctx context.Context, from int, to int) ([]model.Post, error) {
	args := m.Called(ctx, from, to)
	return args.Get(0).([]model.Post), args.Error(1)
}

func NewPost() *PostMock {
	return &PostMock{}
}
