package postrepo

import (
	"context"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"

	"github.com/aedobrynin/soa-hw/posts/internal/repo"
)

type postRepo struct {
	logger  *zap.Logger
	pgxPool *pgxpool.Pool
}

func generatePostId() uuid.UUID {
	return uuid.Must(uuid.NewV4())
}

func (r *postRepo) conn(ctx context.Context) Conn {
	if tx, ok := ctx.Value(repo.CtxKeyTx).(pgx.Tx); ok {
		return tx
	}

	return r.pgxPool
}

func (r *postRepo) WithNewTx(ctx context.Context, f func(ctx context.Context) error) error {
	return r.pgxPool.BeginFunc(ctx, func(tx pgx.Tx) error {
		return f(context.WithValue(ctx, repo.CtxKeyTx, tx))
	})
}

func (r *postRepo) AddPost(ctx context.Context, authorId uuid.UUID, content string) (uuid.UUID, error) {
	postId := generatePostId()

	_, err := r.conn(ctx).Exec(ctx, `INSERT INTO posts.posts (id, author_id, content) VALUES ($1, $2, $3)`,
		postId, authorId, content)
	if err != nil {
		return uuid.Nil, err
	}
	return postId, err
}

func New(logger *zap.Logger, pgxPool *pgxpool.Pool) repo.Post {
	return &postRepo{
		logger:  logger,
		pgxPool: pgxPool,
	}
}
