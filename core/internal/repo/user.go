package repo

import (
	"context"

	"github.com/aedobrynin/soa-hw/core/internal/model"

	"github.com/google/uuid"
)

type AddRequest struct {
	Login    string
	Password string
	Name     *string
	Surname  *string
	Email    *string
	Phone    *string
}

type UpdateRequest struct {
	UserId  uuid.UUID
	Name    *string
	Surname *string
	Email   *string
	Phone   *string
}

type User interface {
	WithNewTx(ctx context.Context, f func(ctx context.Context) error) error
	AddUser(ctx context.Context, request AddRequest) error
	GetUser(ctx context.Context, login string) (*model.User, error)
	ValidateUser(ctx context.Context, login, password string) (*model.User, error)
	UpdateUser(ctx context.Context, request UpdateRequest) error
}
