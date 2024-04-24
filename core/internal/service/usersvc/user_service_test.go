package usersvc_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/aedobrynin/soa-hw/core/internal/repo"
	"github.com/aedobrynin/soa-hw/core/internal/repo/repomock"
	"github.com/aedobrynin/soa-hw/core/internal/service"
	"github.com/aedobrynin/soa-hw/core/internal/service/usersvc"
)

const (
	validLogin    = "valid_login"
	validPassword = "valid_password"
)

// TODO: myb some tests for validation here too?

func TestSignUpHappyPath(t *testing.T) {
	ctx := context.Background()

	const (
		// TODO: fill name, surname, email, phone validation here
		login    = validLogin
		password = validPassword
	)

	userRepo := repomock.NewUser()
	userRepo.On(
		"AddUser",
		mock.MatchedBy(func(ctx context.Context) bool { return true }),
		repo.AddRequest{Login: login, Password: password},
	).Return(nil)

	svc := usersvc.New(userRepo)
	err := svc.SignUp(ctx, service.SignUpRequest{Login: login, Password: password})
	require.Nil(t, err)
}

func TestSignUpLoginIsTaken(t *testing.T) {
	ctx := context.Background()

	const (
		login    = validLogin
		password = validPassword
	)

	userRepo := repomock.NewUser()
	userRepo.On(
		"AddUser",
		mock.MatchedBy(func(ctx context.Context) bool { return true }),
		repo.AddRequest{Login: login, Password: password},
	).Return(nil).Once()

	svc := usersvc.New(userRepo)
	err := svc.SignUp(ctx, service.SignUpRequest{Login: login, Password: password})
	require.Nil(t, err)

	userRepo.On(
		"AddUser",
		mock.MatchedBy(func(ctx context.Context) bool { return true }),
		repo.AddRequest{Login: login, Password: password},
	).Return(repo.ErrLoginTaken).Once()

	err = svc.SignUp(ctx, service.SignUpRequest{Login: login, Password: password})
	require.ErrorIs(t, err, service.ErrLoginTaken)
}

// TODO: test Edit
