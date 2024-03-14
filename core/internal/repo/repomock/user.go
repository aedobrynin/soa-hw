package repomock

import (
	"context"

	"core/internal/model"
	"core/internal/repo"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"
)

var _ repo.User = &UserMock{}

type UserMock struct {
	mock.Mock
}

func (m *UserMock) WithNewTx(ctx context.Context, f func(ctx context.Context) error) error {
	args := m.Called(ctx, f)
	return args.Error(0)
}

func (m *UserMock) AddUser(ctx context.Context, login, password string) error {
	args := m.Called(ctx, login, password)
	return args.Error(0)
}

func (m *UserMock) GetUser(ctx context.Context, login string) (*model.User, error) {
	args := m.Called(ctx, login)
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *UserMock) ValidateUser(ctx context.Context, login, password string) (*model.User, error) {
	args := m.Called(ctx, login, password)
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *UserMock) UpdateName(
	ctx context.Context,
	userId uuid.UUID,
	name string,
) error {
	args := m.Called(ctx, userId, name)
	return args.Error(0)
}

func (m *UserMock) UpdateSurname(
	ctx context.Context,
	userId uuid.UUID,
	surname string,
) error {
	args := m.Called(ctx, userId, surname)
	return args.Error(0)
}

func (m *UserMock) UpdateEmail(
	ctx context.Context,
	userId uuid.UUID,
	email string,
) error {
	args := m.Called(ctx, userId, email)
	return args.Error(0)
}

func (m *UserMock) UpdatePhone(
	ctx context.Context,
	userId uuid.UUID,
	phone string,
) error {
	args := m.Called(ctx, userId, phone)
	return args.Error(0)
}

func NewUser() *UserMock {
	return &UserMock{}
}
