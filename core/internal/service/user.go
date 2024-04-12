package service

import (
	"context"

	"github.com/google/uuid"
)

type User interface {
	SignUp(ctx context.Context, login, password string) error
	ChangeName(ctx context.Context, userId uuid.UUID, name string) error
	ChangeSurname(ctx context.Context, userId uuid.UUID, surname string) error
	ChangeEmail(ctx context.Context, userId uuid.UUID, email string) error
	ChangePhone(ctx context.Context, userId uuid.UUID, phone string) error
}
