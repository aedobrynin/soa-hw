package service

import (
	"context"

	"github.com/gofrs/uuid"
)

type Post interface {
	AddPost(ctx context.Context, authorId uuid.UUID, content string) error
}
