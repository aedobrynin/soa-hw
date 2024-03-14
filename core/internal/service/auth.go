package service

import (
	"context"
	"time"

	"core/internal/model"

	"github.com/gofrs/uuid"
)

type AuthConfig struct {
	SigningKey           string        `yaml:"signing_key"`
	AccessTokenDuration  time.Duration `yaml:"access_token_duration"`
	RefreshTokenDuration time.Duration `yaml:"refresh_token_duration"`
}

type Auth interface {
	Login(ctx context.Context, login, password string) (*model.TokenPair, error)
	ValidateAndRefresh(
		ctx context.Context,
		tokenPair *model.TokenPair,
	) (new *model.TokenPair, userId *uuid.UUID, err error)
}
