package service

import "errors"

var (
	ErrUnauthorized      = errors.New("unauthorized")
	ErrUnsupportedClaims = errors.New("invalid token")
)
