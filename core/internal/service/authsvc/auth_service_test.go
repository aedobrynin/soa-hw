package authsvc_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/aedobrynin/soa-hw/core/internal/model"
	"github.com/aedobrynin/soa-hw/core/internal/repo"
	"github.com/aedobrynin/soa-hw/core/internal/repo/repomock"
	"github.com/aedobrynin/soa-hw/core/internal/service"
	"github.com/aedobrynin/soa-hw/core/internal/service/authsvc"
)

func TestLoginHappyPath(t *testing.T) {
	ctx := context.Background()

	const (
		login    = "login"
		password = "password"
	)

	userRepo := repomock.NewUser()
	userRepo.On(
		"ValidateUser",
		mock.MatchedBy(func(ctx context.Context) bool { return true }),
		login,
		password,
	).Return(&model.User{ID: uuid.New(), Login: login, HashedPassword: []byte("hashed")}, nil).Once()

	config := &service.AuthConfig{
		SigningKey:           "signingKey",
		AccessTokenDuration:  1 * time.Second,
		RefreshTokenDuration: 2 * time.Second,
	}

	svc := authsvc.New(config, userRepo)
	_, err := svc.Login(ctx, login, password)
	require.Nil(t, err)
}

func TestLoginUserNotFound(t *testing.T) {
	ctx := context.Background()

	const (
		login    = "login"
		password = "password"
	)

	var nilUser *model.User = nil

	userRepo := repomock.NewUser()
	userRepo.On(
		"ValidateUser",
		mock.MatchedBy(func(ctx context.Context) bool { return true }),
		login,
		password,
	).Return(nilUser, repo.ErrUserNotFound).Once()

	config := &service.AuthConfig{
		SigningKey:           "signingKey",
		AccessTokenDuration:  1 * time.Second,
		RefreshTokenDuration: 2 * time.Second,
	}

	svc := authsvc.New(config, userRepo)
	_, err := svc.Login(ctx, login, password)
	require.ErrorIs(t, err, service.ErrUserNotFound)
}

func TestLoginWrongPassword(t *testing.T) {
	ctx := context.Background()

	const (
		login    = "login"
		password = "password"
	)

	var nilUser *model.User = nil

	userRepo := repomock.NewUser()
	userRepo.On(
		"ValidateUser",
		mock.MatchedBy(func(ctx context.Context) bool { return true }),
		login,
		password,
	).Return(nilUser, repo.ErrWrongPassword).Once()

	config := &service.AuthConfig{
		SigningKey:           "signingKey",
		AccessTokenDuration:  1 * time.Second,
		RefreshTokenDuration: 2 * time.Second,
	}

	svc := authsvc.New(config, userRepo)
	_, err := svc.Login(ctx, login, password)
	require.ErrorIs(t, err, service.ErrWrongPassword)
}

func TestTokensTTL(t *testing.T) {
	ctx := context.Background()

	userID := uuid.New()

	userRepo := repomock.NewUser()
	userRepo.On(
		"ValidateUser",
		mock.MatchedBy(func(ctx context.Context) bool { return true }),
		"login",
		"password",
	).Return(&model.User{ID: userID, Login: "login", HashedPassword: []byte("hashed")}, nil)

	config := &service.AuthConfig{
		SigningKey:           "signingKey",
		AccessTokenDuration:  1 * time.Second,
		RefreshTokenDuration: 2 * time.Second,
	}

	svc := authsvc.New(config, userRepo)

	initialPair, err := svc.Login(ctx, "login", "password")
	require.Nil(t, err)

	newPair, parsedUserID, err := svc.ValidateAndRefresh(ctx, initialPair)
	require.Nil(t, err)
	require.Equal(t, userID, *parsedUserID)

	require.Equal(t, initialPair.AccessToken, newPair.AccessToken)
	require.Equal(t, initialPair.RefreshToken, newPair.RefreshToken)

	time.Sleep(config.AccessTokenDuration)

	newPair, parsedUserID, err = svc.ValidateAndRefresh(ctx, initialPair)
	require.Nil(t, err)
	require.Equal(t, userID, *parsedUserID)

	require.NotEqual(t, initialPair.AccessToken, newPair.AccessToken)
	require.NotEqual(t, initialPair.RefreshToken, newPair.RefreshToken)

	time.Sleep(config.RefreshTokenDuration)

	_, _, err = svc.ValidateAndRefresh(ctx, newPair)
	require.ErrorIs(t, err, service.ErrUnauthorized)
}
