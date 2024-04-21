package userrepo

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/crypto/bcrypt"

	"github.com/aedobrynin/soa-hw/core/internal/model"
	"github.com/aedobrynin/soa-hw/core/internal/repo"
	"github.com/aedobrynin/soa-hw/core/internal/service"
)

const (
	pgErrorCodeUniqueViolation = "23505"
)

type userRepo struct {
	pgxPool *pgxpool.Pool
}

func generateUserId() uuid.UUID {
	return uuid.New()
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
	request repo.UpdateRequest,
) error {
	if request.Name == nil && request.Surname == nil && request.Email == nil && request.Phone == nil {
		return nil
	}

	var requestBuilder strings.Builder
	requestBuilder.WriteString("UPDATE core.users SET ")
	fields := make([]string, 0)
	args := make([]interface{}, 0)
	if request.Name != nil {
		fields = append(fields, "name")
		args = append(args, *request.Name)
	}
	if request.Surname != nil {
		fields = append(fields, "surname")
		args = append(args, *request.Surname)
	}
	if request.Email != nil {
		fields = append(fields, "email")
		args = append(args, *request.Email)
	}
	if request.Phone != nil {
		fields = append(fields, "phone")
		args = append(args, *request.Phone)
	}
	if len(fields) != len(args) {
		return errors.New("TODO")
	}

	for i := 0; i < len(fields); i++ {
		var queryPart string
		if i+1 != len(fields) {
			queryPart = fmt.Sprintf("%s = $%d, ", fields[i], i+1)
		} else {
			queryPart = fmt.Sprintf("%s = $%d ", fields[i], i+1)
		}
		requestBuilder.WriteString(queryPart)
	}

	args = append(args, request.UserId)
	requestBuilder.WriteString(fmt.Sprintf("WHERE id = $%d", len(args)))

	res, err := r.conn(ctx).Exec(ctx, requestBuilder.String(), args...)
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
