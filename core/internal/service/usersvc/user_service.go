package usersvc

import (
	"context"
	"errors"

	"github.com/aedobrynin/soa-hw/core/internal/repo"
	"github.com/aedobrynin/soa-hw/core/internal/service"
)

type userSvc struct {
	repo repo.User
}

var _ service.User = &userSvc{}

func (s *userSvc) SignUp(ctx context.Context, login, password string) error {
	err := validateLogin(login)
	if err != nil {
		return err
	}
	err = validatePassword(password)
	if err != nil {
		return err
	}
	err = s.repo.AddUser(ctx, login, password)
	if errors.Is(err, repo.ErrLoginTaken) {
		return service.ErrLoginTaken
	}
	return err
}

func (s *userSvc) Edit(ctx context.Context, request service.EditRequest) error {
	if request.Name != nil {
		if err := validateName(*request.Name); err != nil {
			return err
		}
	}
	if request.Surname != nil {
		if err := validateSurname(*request.Surname); err != nil {
			return err
		}
	}
	if request.Email != nil {
		if err := validateEmail(*request.Email); err != nil {
			return err
		}
	}
	if request.Phone != nil {
		if err := validatePhone(*request.Phone); err != nil {
			return err
		}
	}

	err := s.repo.UpdateUser(
		ctx,
		repo.UpdateRequest{
			UserId:  request.UserId,
			Name:    request.Name,
			Surname: request.Surname,
			Email:   request.Email,
			Phone:   request.Phone,
		},
	)
	switch {
	case errors.Is(err, repo.ErrUserNotFound):
		return service.ErrUserNotFound
	case err != nil:
		return err
	}
	return nil
}

func New(repo repo.User) service.User {
	return &userSvc{repo: repo}
}
