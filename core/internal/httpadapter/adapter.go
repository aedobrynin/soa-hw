package httpadapter

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"

	"core/internal/httpadapter/codegen"
	"core/internal/logger"
	"core/internal/model"
	"core/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/juju/zaputil/zapctx"
	"go.uber.org/zap"
)

type adapter struct {
	cfg *Config

	authService service.Auth
	userService service.User

	server *http.Server

	logger *zap.Logger
}

// (POST /v1/sign_up)
func (a *adapter) PostV1SignUp(
	ctx context.Context,
	request codegen.PostV1SignUpRequestObject,
) (codegen.PostV1SignUpResponseObject, error) {
	err := a.userService.SignUp(ctx, request.Body.Login, request.Body.Password)
	switch {
	case errors.Is(err, service.ErrLoginValidation) || errors.Is(err, service.ErrPasswordValidation) || errors.Is(err, service.ErrLoginTaken):
		return codegen.PostV1SignUp422JSONResponse(codegen.ErrorMessage{Error: err.Error()}), nil
	case err != nil:
		return nil, err
	default:
		return codegen.PostV1SignUp200Response{}, nil
	}
}

// (POST /v1/auth)
func (a *adapter) PostV1Auth(
	ctx context.Context,
	request codegen.PostV1AuthRequestObject,
) (codegen.PostV1AuthResponseObject, error) {
	tokenPair, err := a.authService.Login(ctx, request.Body.Login, request.Body.Password)
	// TODO: use refresh token too
	switch {
	case errors.Is(err, service.ErrUserNotFound) || errors.Is(err, service.ErrWrongPassword):
		return codegen.PostV1Auth401JSONResponse(codegen.ErrorMessage{Error: err.Error()}), nil
	case err != nil:
		return nil, err
	}
	return codegen.PostV1Auth200Response{
		Headers: codegen.PostV1Auth200ResponseHeaders{
			SetCookie: fmt.Sprintf(
				"%s=%s",
				a.cfg.AccessTokenCookie,
				tokenPair.AccessToken,
			),
		},
	}, nil
}

// (POST /v1/change_name)
func (a *adapter) PostV1ChangeName(
	ctx context.Context,
	request codegen.PostV1ChangeNameRequestObject,
) (codegen.PostV1ChangeNameResponseObject, error) {
	// TODO: use refresh token too
	_, userId, err := a.authService.ValidateAndRefresh(
		ctx,
		&model.TokenPair{AccessToken: request.Params.XSESSION, RefreshToken: ""},
	)
	switch {
	case errors.Is(err, service.ErrUnsupportedClaims) || errors.Is(err, service.ErrUnauthorized):
		return codegen.PostV1ChangeName401Response{}, nil
	case err != nil:
		return nil, err
	}

	err = a.userService.ChangeName(ctx, *userId, request.Body.Name)
	switch {
	case errors.Is(err, service.ErrNameValidation):
		return codegen.PostV1ChangeName422JSONResponse(codegen.ErrorMessage{Error: err.Error()}), nil
	case errors.Is(err, service.ErrUserNotFound):
		return codegen.PostV1ChangeName401Response{}, nil
	case err != nil:
		return nil, err
	}
	return codegen.PostV1ChangeName200Response{}, nil
}

// (POST /v1/change_surname)
func (a *adapter) PostV1ChangeSurname(
	ctx context.Context,
	request codegen.PostV1ChangeSurnameRequestObject,
) (codegen.PostV1ChangeSurnameResponseObject, error) {
	// TODO: use refresh token too
	_, userId, err := a.authService.ValidateAndRefresh(
		ctx,
		&model.TokenPair{AccessToken: request.Params.XSESSION, RefreshToken: ""},
	)
	switch {
	case errors.Is(err, service.ErrUnsupportedClaims) || errors.Is(err, service.ErrUnauthorized):
		return codegen.PostV1ChangeSurname401Response{}, nil
	case err != nil:
		return nil, err
	}

	err = a.userService.ChangeSurname(ctx, *userId, request.Body.Surname)
	switch {
	case errors.Is(err, service.ErrSurnameValidation):
		return codegen.PostV1ChangeSurname422JSONResponse(codegen.ErrorMessage{Error: err.Error()}), nil
	case errors.Is(err, service.ErrUserNotFound):
		return codegen.PostV1ChangeSurname401Response{}, nil
	case err != nil:
		return nil, err
	}
	return codegen.PostV1ChangeSurname200Response{}, nil
}

// (POST /v1/change_email)
func (a *adapter) PostV1ChangeEmail(
	ctx context.Context,
	request codegen.PostV1ChangeEmailRequestObject,
) (codegen.PostV1ChangeEmailResponseObject, error) {
	// TODO: use refresh token too
	_, userId, err := a.authService.ValidateAndRefresh(
		ctx,
		&model.TokenPair{AccessToken: request.Params.XSESSION, RefreshToken: ""},
	)
	switch {
	case errors.Is(err, service.ErrUnsupportedClaims) || errors.Is(err, service.ErrUnauthorized):
		return codegen.PostV1ChangeEmail401Response{}, nil
	case err != nil:
		return nil, err
	}

	err = a.userService.ChangeEmail(ctx, *userId, request.Body.Email)
	switch {
	case errors.Is(err, service.ErrEmailValidation):
		return codegen.PostV1ChangeEmail422JSONResponse(codegen.ErrorMessage{Error: err.Error()}), nil
	case errors.Is(err, service.ErrUserNotFound):
		return codegen.PostV1ChangeEmail401Response{}, nil
	case err != nil:
		return nil, err
	}
	return codegen.PostV1ChangeEmail200Response{}, nil

}

// (POST /v1/change_phone)
func (a *adapter) PostV1ChangePhone(
	ctx context.Context,
	request codegen.PostV1ChangePhoneRequestObject,
) (codegen.PostV1ChangePhoneResponseObject, error) {
	// TODO: use refresh token too
	_, userId, err := a.authService.ValidateAndRefresh(
		ctx,
		&model.TokenPair{AccessToken: request.Params.XSESSION, RefreshToken: ""},
	)
	switch {
	case errors.Is(err, service.ErrUnsupportedClaims) || errors.Is(err, service.ErrUnauthorized):
		return codegen.PostV1ChangePhone401Response{}, nil
	case err != nil:
		return nil, err
	}

	err = a.userService.ChangePhone(ctx, *userId, request.Body.Phone)
	switch {
	case errors.Is(err, service.ErrPhoneValidation):
		return codegen.PostV1ChangePhone422JSONResponse(codegen.ErrorMessage{Error: err.Error()}), nil
	case errors.Is(err, service.ErrUserNotFound):
		return codegen.PostV1ChangePhone401Response{}, nil
	case err != nil:
		return nil, err
	}
	return codegen.PostV1ChangePhone200Response{}, nil
}

func (a *adapter) Serve() error {
	logger, err := logger.GetLogger(true)
	if err != nil {
		return err
	}
	a.logger = logger

	handler_opts := codegen.StrictHTTPServerOptions{
		RequestErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		},
		ResponseErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
			logger.Error(err.Error())
			http.Error(w, "server got itself in trouble", http.StatusInternalServerError)
		},
	}
	strict_handler := codegen.NewStrictHandlerWithOptions(a, make([]codegen.StrictMiddlewareFunc, 0), handler_opts)
	options := codegen.ChiServerOptions{
		BaseURL:     a.cfg.BasePath,
		BaseRouter:  chi.NewRouter(),
		Middlewares: make([]codegen.MiddlewareFunc, 0),
	}

	a.server = &http.Server{
		Handler: codegen.HandlerWithOptions(strict_handler, options),
		Addr:    a.cfg.ServeAddress,
		BaseContext: func(_ net.Listener) context.Context {
			return zapctx.WithLogger(context.Background(), logger)
		},
	}
	if a.cfg.UseTLS {
		return a.server.ListenAndServeTLS(a.cfg.TLSCrtFile, a.cfg.TLSKeyFile)
	}
	logger.Sugar().Infof("Server started on addr: %s", a.server.Addr)
	return a.server.ListenAndServe()
}

func (a *adapter) Shutdown(ctx context.Context) {
	logger := zapctx.Logger(ctx)
	logger.Info("Server is shutting down...")
	_ = a.server.Shutdown(ctx)
	logger.Info("Server is closed")
}

func NewAdapter(
	config *Config,
	authService service.Auth,
	userService service.User) Adapter {
	return &adapter{
		cfg:         config,
		authService: authService,
		userService: userService,
	}
}
