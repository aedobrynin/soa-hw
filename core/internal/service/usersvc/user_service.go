package usersvc

import (
	"context"
	"errors"

	"github.com/aedobrynin/soa-hw/core/internal/repo"
	"github.com/aedobrynin/soa-hw/core/internal/service"

	"github.com/gofrs/uuid"
)

type userSvc struct {
	repo repo.User
}

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

func (s *userSvc) ChangeName(ctx context.Context, userId uuid.UUID, name string) error {
	if err := validateName(name); err != nil {
		return err
	}

	err := s.repo.UpdateName(ctx, userId, name)
	switch {
	case errors.Is(err, repo.ErrUserNotFound):
		return service.ErrUserNotFound
	case err != nil:
		return err
	}
	return nil
}

func (s *userSvc) ChangeSurname(ctx context.Context, userId uuid.UUID, surname string) error {
	if err := validateSurname(surname); err != nil {
		return err
	}

	err := s.repo.UpdateSurname(ctx, userId, surname)
	switch {
	case errors.Is(err, repo.ErrUserNotFound):
		return service.ErrUserNotFound
	case err != nil:
		return err
	}
	return nil
}

func (s *userSvc) ChangeEmail(ctx context.Context, userId uuid.UUID, email string) error {
	if err := validateEmail(email); err != nil {
		return err
	}

	err := s.repo.UpdateEmail(ctx, userId, email)
	switch {
	case errors.Is(err, repo.ErrUserNotFound):
		return service.ErrUserNotFound
	case err != nil:
		return err
	}
	return nil
}

func (s *userSvc) ChangePhone(ctx context.Context, userId uuid.UUID, phone string) error {
	if err := validatePhone(phone); err != nil {
		return err
	}

	err := s.repo.UpdatePhone(ctx, userId, phone)
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
