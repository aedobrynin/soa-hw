package userrepo

import (
	"context"
	"errors"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/crypto/bcrypt"

	"core/internal/model"
	"core/internal/repo"
	"core/internal/service"
)

const (
	pgErrorCodeUniqueViolation = "23505"
)

type userRepo struct {
	pgxPool *pgxpool.Pool
}

func generateUserId() uuid.UUID {
	return uuid.Must(uuid.NewV4())
}

func (r *userRepo) conn(ctx context.Context) Conn {
	if tx, ok := ctx.Value(repo.CtxKeyTx).(pgx.Tx); ok {
		return tx
	}

	return r.pgxPool
}

func (r *userRepo) WithNewTx(ctx context.Context, f func(ctx context.Context) error) error {
	return r.pgxPool.BeginFunc(ctx, func(tx pgx.Tx) error {
		return f(context.WithValue(ctx, repo.CtxKeyTx, tx))
	})
}

func (r *userRepo) AddUser(ctx context.Context, login, password string) error {
	password_hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	userId := generateUserId()

	_, err = r.conn(ctx).Exec(ctx, `INSERT INTO core.users (id, login, password_hash) VALUES ($1, $2, $3)`,
		userId, login, password_hash)
	var pgxErr *pgconn.PgError
	if errors.As(err, &pgxErr) {
		if pgxErr.Code == pgErrorCodeUniqueViolation {
			return repo.ErrLoginTaken
		}
	}
	if err != nil {
		return err
	}

	return nil
}

func (r *userRepo) GetUser(ctx context.Context, login string) (*model.User, error) {
	var user model.User

	row := r.conn(ctx).
		QueryRow(ctx, `SELECT id, login, password_hash, name, surname, email, phone FROM core.users WHERE login = $1`, login)
	if err := row.Scan(&user.Id, &user.Login, &user.HashedPassword, &user.Name, &user.Surname, &user.Email, &user.Phone); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repo.ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (r *userRepo) ValidateUser(ctx context.Context, login, password string) (*model.User, error) {
	user, err := r.GetUser(ctx, login)
	switch {
	case errors.Is(err, repo.ErrUserNotFound):
		return nil, repo.ErrUserNotFound
	case err != nil:
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword(user.HashedPassword, []byte(password)); err != nil {
		return nil, repo.ErrWrongPassword
	}

	return user, nil
}

func (r *userRepo) UpdateUser(
	ctx context.Context,
	userId uuid.UUID,
	name *string,
	surname *string,
	email *string,
	phone *string,
) error {
	res, err := r.conn(ctx).Exec(ctx, `UPDATE core.users SET name = $1 WHERE id = $2`, *name, userId)
	if err != nil {
		return err
	}
	if res.RowsAffected() != 1 {
		return repo.ErrUserNotFound
	}
	return nil
}

func (r *userRepo) UpdateName(ctx context.Context, userId uuid.UUID, name string) error {
	res, err := r.conn(ctx).Exec(ctx, `UPDATE core.users SET name = $1 WHERE id = $2`, name, userId)
	if err != nil {
		return err
	}
	if res.RowsAffected() != 1 {
		return repo.ErrUserNotFound
	}
	return nil
}

func (r *userRepo) UpdateSurname(ctx context.Context, userId uuid.UUID, surname string) error {
	res, err := r.conn(ctx).Exec(ctx, `UPDATE core.users SET surname = $1 WHERE id = $2`, surname, userId)
	if err != nil {
		return err
	}
	if res.RowsAffected() != 1 {
		return repo.ErrUserNotFound
	}
	return nil
}
func (r *userRepo) UpdateEmail(ctx context.Context, userId uuid.UUID, email string) error {
	res, err := r.conn(ctx).Exec(ctx, `UPDATE core.users SET email = $1 WHERE id = $2`, email, userId)
	if err != nil {
		return err
	}
	if res.RowsAffected() != 1 {
		return repo.ErrUserNotFound
	}
	return nil
}
func (r *userRepo) UpdatePhone(ctx context.Context, userId uuid.UUID, phone string) error {
	res, err := r.conn(ctx).Exec(ctx, `UPDATE core.users SET phone = $1 WHERE id = $2`, phone, userId)
	if err != nil {
		return err
	}
	if res.RowsAffected() != 1 {
		return repo.ErrUserNotFound
	}
	return nil
}

func New(config *service.AuthConfig, pgxPool *pgxpool.Pool) repo.User {
	return &userRepo{
		pgxPool: pgxPool,
	}
}
