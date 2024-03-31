package repomock

import (
	"context"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"
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

func NewPost() *PostMock {
	return &PostMock{}
}
