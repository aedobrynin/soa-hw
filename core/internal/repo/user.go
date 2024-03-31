package repo

import (
	"context"

	"github.com/aedobrynin/soa-hw/core/internal/model"

	"github.com/gofrs/uuid"
)

type User interface {
	WithNewTx(ctx context.Context, f func(ctx context.Context) error) error
	AddUser(ctx context.Context, login, password string) error
	GetUser(ctx context.Context, login string) (*model.User, error)
	ValidateUser(ctx context.Context, login, password string) (*model.User, error)

	UpdateName(ctx context.Context, userId uuid.UUID, name string) error
	UpdateSurname(ctx context.Context, userId uuid.UUID, surname string) error
	UpdateEmail(ctx context.Context, userId uuid.UUID, email string) error
	UpdatePhone(ctx context.Context, userId uuid.UUID, phone string) error
}
