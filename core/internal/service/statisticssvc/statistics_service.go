package statisticssvc

import (
	"context"
	"fmt"

	"github.com/aedobrynin/soa-hw/core/internal/model"
	"github.com/aedobrynin/soa-hw/core/internal/repo"
	"github.com/aedobrynin/soa-hw/core/internal/service"
)

type statisticsSvc struct {
	repo repo.Statistics
}

var _ service.Statistics = &statisticsSvc{}

func (s *statisticsSvc) AccountPostView(ctx context.Context, view model.PostView) error {
	err := s.repo.PushPostView(ctx, view)
	if err != nil {
		return fmt.Errorf("error on accounting post view: %v", err)
	}
	return nil
}

func (s *statisticsSvc) AccountPostLike(ctx context.Context, like model.PostLike) error {
	err := s.repo.PushPostLike(ctx, like)
	if err != nil {
		return fmt.Errorf("error on accounting post like: %v", err)
	}
	return nil
}

func New(repo repo.Statistics) service.Statistics {
	return &statisticsSvc{repo: repo}
}
