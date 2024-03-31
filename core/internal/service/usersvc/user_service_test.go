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

func TestSignUpLoginValidation(t *testing.T) {
	ctx := context.Background()

	userRepo := repomock.NewUser()
	svc := usersvc.New(userRepo)

	// Too short
	err := svc.SignUp(ctx, "", "good")
	require.ErrorIs(t, err, service.ErrLoginValidation)

	// Too long
	err = svc.SignUp(ctx, strings.Repeat("a", 26), "good")
	require.ErrorIs(t, err, service.ErrLoginValidation)

	// Unsupported letter
	err = svc.SignUp(ctx, "русские_буквы", "good")
	require.ErrorIs(t, err, service.ErrLoginValidation)

	// Symbol
	err = svc.SignUp(ctx, "okokok@", "good")
	require.ErrorIs(t, err, service.ErrLoginValidation)
}

func TestSignUpPasswordValidation(t *testing.T) {
	ctx := context.Background()

	userRepo := repomock.NewUser()
	svc := usersvc.New(userRepo)

	// Too short
	err := svc.SignUp(ctx, "good_login", strings.Repeat("a", 7))
	require.ErrorIs(t, err, service.ErrPasswordValidation)

	// Too long
	err = svc.SignUp(ctx, "good_login", strings.Repeat("a", 256))
	require.ErrorIs(t, err, service.ErrPasswordValidation)

	// Unsupported letter
	err = svc.SignUp(ctx, "good_login", "русские_буквы")
	require.ErrorIs(t, err, service.ErrPasswordValidation)
}

func TestSignUpHappyPath(t *testing.T) {
	ctx := context.Background()

	const (
		login          = "login"
		password       = "passworddddd"
		hashedPassword = "hashedPassword"
	)

	userRepo := repomock.NewUser()
	userRepo.On(
		"AddUser",
		mock.MatchedBy(func(ctx context.Context) bool { return true }),
		login,
		password,
	).Return(nil)

	svc := usersvc.New(userRepo)
	err := svc.SignUp(ctx, login, password)
	require.Nil(t, err)
}

func TestSignUpLoginIsTaken(t *testing.T) {
	ctx := context.Background()

	const (
		login          = "login"
		password       = "passworddddd"
		hashedPassword = "hashedPassword"
	)

	userRepo := repomock.NewUser()
	userRepo.On(
		"AddUser",
		mock.MatchedBy(func(ctx context.Context) bool { return true }),
		login,
		password,
	).Return(nil).Once()

	svc := usersvc.New(userRepo)
	err := svc.SignUp(ctx, login, password)
	require.Nil(t, err)

	userRepo.On(
		"AddUser",
		mock.MatchedBy(func(ctx context.Context) bool { return true }),
		login,
		password,
	).Return(repo.ErrLoginTaken).Once()

	err = svc.SignUp(ctx, login, password)
	require.ErrorIs(t, err, service.ErrLoginTaken)
}

// TODO: test ChangeName
// TODO: test ChangeSurname
// TODO: test ChangeEmail
// TODO: test ChangePhone
