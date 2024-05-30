package service

import (
	"context"

	"github.com/aedobrynin/soa-hw/core/internal/model"
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
	UserID  model.UserID
	Name    *string
	Surname *string
	Email   *string
	Phone   *string
}

type User interface {
	SignUp(ctx context.Context, request SignUpRequest) error
	Edit(ctx context.Context, request EditRequest) error
	GetUser(ctx context.Context, userID model.UserID) (*model.User, error)
}
