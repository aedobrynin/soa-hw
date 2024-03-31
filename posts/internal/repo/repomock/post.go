package repomock

import (
	"context"

	"github.com/google/uuid"
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

func (m *PostMock) AddPost(ctx context.Context, authorId uuid.UUID, content string) (uuid.UUID, error) {
	args := m.Called(ctx, authorId, content)
	return args.Get(0).(uuid.UUID), args.Error(1)
}

func (m *PostMock) GetPost(ctx context.Context, postId uuid.UUID) (*model.Post, error) {
	args := m.Called(ctx, postId)
	return args.Get(0).(*model.Post), args.Error(1)
}

func (m *PostMock) EditPost(ctx context.Context, postId uuid.UUID, content string) error {
	args := m.Called(ctx, postId, content)
	return args.Error(0)
}

func (m *PostMock) DeletePost(ctx context.Context, postId uuid.UUID) error {
	args := m.Called(ctx, postId)
	return args.Error(0)
}

func NewPost() *PostMock {
	return &PostMock{}
}
