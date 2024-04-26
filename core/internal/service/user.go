package service

import (
	"context"

	"github.com/google/uuid"
)

type SignUpRequest struct {
	Login    string
	Password string
	Name     *string
	Surname  *string
	Email    *string
	Phone    *string
}

type EditRequest struct {
	UserID  uuid.UUID
	Name    *string
	Surname *string
	Email   *string
	Phone   *string
}

type User interface {
	SignUp(ctx context.Context, request SignUpRequest) error
	Edit(ctx context.Context, request EditRequest) error
}
