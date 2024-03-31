package service

import "errors"

var (
	ErrContentIsEmpty          = errors.New("content is empty")
	ErrPostNotFound            = errors.New("post not found")
	ErrInsufficientPermissions = errors.New("insufficient permissions")
)
