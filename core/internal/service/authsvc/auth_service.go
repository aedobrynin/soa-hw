package authsvc

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aedobrynin/soa-hw/core/internal/model"
	"github.com/aedobrynin/soa-hw/core/internal/repo"
	"github.com/aedobrynin/soa-hw/core/internal/service"

	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	jwt.RegisteredClaims

	UserID model.UserID `json:"user_id"`
}

type authService struct {
	repo repo.User

	signingKey           string
	accessTokenDuration  time.Duration
	refreshTokenDuration time.Duration
}

func (s *authService) makeToken(userID model.UserID, duration time.Duration) (string, error) {
	now := time.Now().UTC()

	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(duration)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
		UserID: userID,
	}

	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(s.signingKey))
}

func (s *authService) newTokenPair(userID model.UserID) (*model.TokenPair, error) {
	accessToken, err := s.makeToken(userID, s.accessTokenDuration)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.makeToken(userID, s.refreshTokenDuration)
	if err != nil {
		return nil, err
	}

	return &model.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *authService) parseTokenString(tokenString string) (*jwt.Token, error) {
	return jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Method.Alg())
		}

		return []byte(s.signingKey), nil
	})
}

func (s *authService) Login(ctx context.Context, login, password string) (*model.TokenPair, error) {
	user, err := s.repo.ValidateUser(ctx, login, password)
	switch {
	case errors.Is(err, repo.ErrUserNotFound):
		return nil, service.ErrUserNotFound
	case errors.Is(err, repo.ErrWrongPassword):
		return nil, service.ErrWrongPassword
	case err != nil:
		return nil, err
	default:
		return s.newTokenPair(user.ID)
	}
}

func (s *authService) ValidateAndRefresh(
	ctx context.Context,
	tokenPair *model.TokenPair,
) (new *model.TokenPair, userID *model.UserID, err error) {
	accessToken, err := s.parseTokenString(tokenPair.AccessToken)

	switch v := err.(type) {
	case nil:
		claims, ok := accessToken.Claims.(*Claims)
		if !ok {
			return nil, nil, service.ErrUnsupportedClaims
		}
		return tokenPair, &claims.UserID, nil

	case *jwt.ValidationError:
		if v.Errors&jwt.ValidationErrorExpired == 0 {
			return nil, nil, err
		}

		_, err = s.parseTokenString(tokenPair.RefreshToken)
		if err != nil {
			return nil, nil, fmt.Errorf("%w: refresh token not valid: %s", service.ErrUnauthorized, err)
		}

		claims, ok := accessToken.Claims.(*Claims)
		if !ok {
			return nil, nil, service.ErrUnsupportedClaims
		}

		var newTokenPair *model.TokenPair
		newTokenPair, err = s.newTokenPair(claims.UserID)
		return newTokenPair, &claims.UserID, err
	}

	return nil, nil, err
}

func New(config *service.AuthConfig, repo repo.User) service.Auth {
	return &authService{
		repo:                 repo,
		signingKey:           config.SigningKey,
		accessTokenDuration:  config.AccessTokenDuration,
		refreshTokenDuration: config.RefreshTokenDuration,
	}
}
