package repo

import "errors"

var (
	ErrLoginTaken    = errors.New("login is taken")
	ErrUserNotFound  = errors.New("user not found")
	ErrWrongPassword = errors.New("wrong password")
)
