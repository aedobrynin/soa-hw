package repo

import (
	"context"

	"github.com/aedobrynin/soa-hw/core/internal/model"
)

type Statistics interface {
	PushPostView(ctx context.Context, view model.PostView) error
	PushPostLike(ctx context.Context, like model.PostLike) error
	Stop(ctx context.Context)
}
