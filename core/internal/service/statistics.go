package service

import (
	"context"

	"github.com/aedobrynin/soa-hw/core/internal/model"
)

type Statistics interface {
	AccountPostView(ctx context.Context, view model.PostView) error
	AccountPostLike(ctx context.Context, like model.PostLike) error
}
