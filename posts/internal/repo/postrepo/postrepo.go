package postrepo

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"

	"github.com/aedobrynin/soa-hw/posts/internal/model"
	"github.com/aedobrynin/soa-hw/posts/internal/repo"
)

type postRepo struct {
	logger  *zap.Logger
	pgxPool *pgxpool.Pool
}

func generatePostId() uuid.UUID {
	return uuid.New()
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

func (r *postRepo) GetPost(ctx context.Context, postId uuid.UUID) (*model.Post, error) {
	var post model.Post

	row := r.conn(ctx).
		QueryRow(ctx, `SELECT id, author_id, content, created_ts, updated_ts FROM posts.posts WHERE id = $1`, postId)
	if err := row.Scan(&post.Id, &post.AuthorId, &post.Content, &post.CreatedTs, &post.UpdatedTs); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repo.ErrPostNotFound
		}
		return nil, err
	}

	return &post, nil
}

func (r *postRepo) EditPost(ctx context.Context, postId uuid.UUID, content string) error {
	res, err := r.conn(ctx).Exec(ctx, `UPDATE posts.posts SET content = $1 WHERE id = $2`, content, postId)
	if err != nil {
		return err
	}
	if res.RowsAffected() != 1 {
		return repo.ErrPostNotFound
	}
	return nil
}

func (r *postRepo) DeletePost(ctx context.Context, postId uuid.UUID) error {
	res, err := r.conn(ctx).Exec(ctx, `DELETE FROM posts.posts WHERE id = $1`, postId)
	if err != nil {
		return err
	}
	if res.RowsAffected() != 1 {
		return repo.ErrPostNotFound
	}
	return nil
}

func (r *postRepo) ListPosts(ctx context.Context, from int, to int) ([]model.Post, error) {
	// TODO: better
	defer r.logger.Sync()
	r.logger.Debug("ListPosts: executing query")
	rows, err := r.conn(ctx).
		Query(ctx, `SELECT id, author_id, content, created_ts, updated_ts FROM posts.posts ORDER BY created_ts DESC LIMIT $1 OFFSET $2`, to-from, from)
	if err != nil {
		r.logger.Sugar().Debugf("error %v", err)
		return nil, err
	}

	posts := make([]model.Post, 0)
	for rows.Next() {
		var post model.Post
		if err := rows.Scan(&post.Id, &post.AuthorId, &post.Content, &post.CreatedTs, &post.UpdatedTs); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				r.logger.Debug("ListPosts: return 0 posts")
				return posts, nil
			}
			r.logger.Sugar().Debugf("ListPosts: error %v", err)
			return nil, err
		}
		posts = append(posts, post)
	}
	r.logger.Sugar().Debugf("ListPosts: return %d posts", len(posts))
	return posts, nil
}

func New(logger *zap.Logger, pgxPool *pgxpool.Pool) repo.Post {
	return &postRepo{
		logger:  logger,
		pgxPool: pgxPool,
	}
}
