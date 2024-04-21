package usersvc_test

import (
	"context"
	"strings"
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

func TestSignUpLoginValidation(t *testing.T) {
	ctx := context.Background()

	userRepo := repomock.NewUser()
	svc := usersvc.New(userRepo)

	// Too short
	err := svc.SignUp(ctx, service.SignUpRequest{Login: "", Password: validPassword})
	require.ErrorIs(t, err, service.ErrLoginValidation)

	// Too long
	err = svc.SignUp(ctx, service.SignUpRequest{Login: strings.Repeat("a", 26), Password: validPassword})
	require.ErrorIs(t, err, service.ErrLoginValidation)

	// Unsupported letter
	err = svc.SignUp(ctx, service.SignUpRequest{Login: "русские_буквы", Password: validPassword})
	require.ErrorIs(t, err, service.ErrLoginValidation)

	// Symbol
	err = svc.SignUp(ctx, service.SignUpRequest{Login: "okokok@", Password: validPassword})
	require.ErrorIs(t, err, service.ErrLoginValidation)
}

func TestSignUpPasswordValidation(t *testing.T) {
	ctx := context.Background()

	userRepo := repomock.NewUser()
	svc := usersvc.New(userRepo)

	// Too short
	err := svc.SignUp(ctx, service.SignUpRequest{Login: validLogin, Password: strings.Repeat("a", 7)})
	require.ErrorIs(t, err, service.ErrPasswordValidation)

	// Too long
	err = svc.SignUp(ctx, service.SignUpRequest{Login: validLogin, Password: strings.Repeat("a", 256)})
	require.ErrorIs(t, err, service.ErrPasswordValidation)

	// Unsupported letter
	err = svc.SignUp(ctx, service.SignUpRequest{Login: validLogin, Password: "русские_буквы"})
	require.ErrorIs(t, err, service.ErrPasswordValidation)
}

// TODO test name, surname, email, phone validation

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
