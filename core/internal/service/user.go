package service

import (
	"context"

	"github.com/google/uuid"
)

type EditRequest struct {
	UserId  uuid.UUID
	Name    *string
	Surname *string
	Email   *string
	Phone   *string
}

type User interface {
	SignUp(ctx context.Context, login, password string) error
	Edit(ctx context.Context, request EditRequest) error
}
