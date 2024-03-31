package repo

import (
	"context"

	"github.com/gofrs/uuid"
)

type Post interface {
	WithNewTx(ctx context.Context, f func(ctx context.Context) error) error
	AddPost(ctx context.Context, authorId uuid.UUID, content string) (postId uuid.UUID, err error)
}
