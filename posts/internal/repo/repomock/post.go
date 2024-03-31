package repomock

import (
	"context"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"
)

type UserMock struct {
	mock.Mock
}

func (m *UserMock) WithNewTx(ctx context.Context, f func(ctx context.Context) error) error {
	args := m.Called(ctx, f)
	return args.Error(0)
}

func (m *UserMock) AddPost(ctx context.Context, authorId uuid.UUID, content string) error {
	args := m.Called(ctx, authorId, content)
	return args.Error(0)
}

func NewUser() *UserMock {
	return &UserMock{}
}
