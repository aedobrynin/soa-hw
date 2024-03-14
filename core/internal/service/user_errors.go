package service

import "errors"

var (
	ErrLoginValidation    = errors.New("bad login")
	ErrPasswordValidation = errors.New("bad password")
	ErrNameValidation     = errors.New("bad name")
	ErrSurnameValidation  = errors.New("bad surname")
	ErrEmailValidation    = errors.New("bad email")
	ErrPhoneValidation    = errors.New("bad phone")
	ErrWrongPassword      = errors.New("wrong password")
	ErrLoginTaken         = errors.New("login is taken")
	ErrUserNotFound       = errors.New("user not found")
)
