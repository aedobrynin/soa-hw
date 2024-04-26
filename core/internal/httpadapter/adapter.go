package httpadapter

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"

	"github.com/aedobrynin/soa-hw/core/internal/clients"
	"github.com/aedobrynin/soa-hw/core/internal/httpadapter/codegen"
	"github.com/aedobrynin/soa-hw/core/internal/logger"
	"github.com/aedobrynin/soa-hw/core/internal/model"
	"github.com/aedobrynin/soa-hw/core/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/juju/zaputil/zapctx"
	"go.uber.org/zap"
)

type adapter struct {
	cfg *Config

	authService service.Auth
	userService service.User

	postsClient clients.PostsClient

	server *http.Server

	logger *zap.Logger
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

// (POST /v1/users)
func (a *adapter) PostV1Users(
	ctx context.Context,
	request codegen.PostV1UsersRequestObject,
) (codegen.PostV1UsersResponseObject, error) {
	err := a.userService.SignUp(
		ctx,
		service.SignUpRequest{
			Login:    request.Body.Login,
			Password: request.Body.Password,
			Name:     request.Body.Name,
			Surname:  request.Body.Surname,
			Email:    request.Body.Email,
			Phone:    request.Body.Phone,
		},
	)
	switch {
	case errors.Is(err, service.ErrLoginValidation) ||
		errors.Is(err, service.ErrPasswordValidation) ||
		errors.Is(err, service.ErrNameValidation) ||
		errors.Is(err, service.ErrSurnameValidation) ||
		errors.Is(err, service.ErrEmailValidation) ||
		errors.Is(err, service.ErrPhoneValidation) ||
		errors.Is(err, service.ErrLoginTaken):
		return codegen.PostV1Users422JSONResponse(codegen.ErrorMessage{Error: err.Error()}), nil
	case err != nil:
		return nil, err
	default:
		return codegen.PostV1Users200Response{}, nil
	}
}

// (PATCH /v1/users/{user_id})
func (a *adapter) PatchV1UsersUserId(
	ctx context.Context,
	request codegen.PatchV1UsersUserIdRequestObject,
) (codegen.PatchV1UsersUserIdResponseObject, error) {
	// TODO: use refresh token too
	_, userID, err := a.authService.ValidateAndRefresh(
		ctx,
		&model.TokenPair{AccessToken: request.Params.XSESSION, RefreshToken: ""},
	)
	switch {
	case errors.Is(err, service.ErrUnsupportedClaims) || errors.Is(err, service.ErrUnauthorized) || request.UserId != userID.String():
		return codegen.PatchV1UsersUserId401Response{}, nil
	case err != nil:
		return nil, err
	}

	err = a.userService.Edit(ctx, service.EditRequest{
		UserID:  *userID,
		Name:    request.Body.Name,
		Surname: request.Body.Surname,
		Email:   request.Body.Email,
		Phone:   request.Body.Phone,
	})
	switch {
	case errors.Is(err, service.ErrNameValidation) ||
		errors.Is(err, service.ErrSurnameValidation) ||
		errors.Is(err, service.ErrEmailValidation) ||
		errors.Is(err, service.ErrPhoneValidation):
		return codegen.PatchV1UsersUserId422JSONResponse(codegen.ErrorMessage{Error: err.Error()}), nil
	case errors.Is(err, service.ErrUserNotFound):
		return codegen.PatchV1UsersUserId401Response{}, nil
	case err != nil:
		return nil, err
	}
	return codegen.PatchV1UsersUserId200Response{}, nil
}

// (POST /v1/posts)
func (a *adapter) PostV1Posts(
	ctx context.Context,
	request codegen.PostV1PostsRequestObject,
) (codegen.PostV1PostsResponseObject, error) {
	// TODO: use refresh token too
	// TODO: make it helper function
	_, userID, err := a.authService.ValidateAndRefresh(
		ctx,
		&model.TokenPair{AccessToken: request.Params.XSESSION, RefreshToken: ""},
	)
	switch {
	case errors.Is(err, service.ErrUnsupportedClaims) || errors.Is(err, service.ErrUnauthorized):
		return codegen.PostV1Posts401Response{}, nil
	case err != nil:
		return nil, err
	}

	_, err = a.postsClient.CreatePost(ctx, *userID, request.Body.Content)
	if errors.Is(err, clients.ErrContentIsEmpty) {
		return codegen.PostV1Posts422JSONResponse(codegen.ErrorMessage{Error: err.Error()}), nil
	}
	if err != nil {
		return nil, err
	}
	return codegen.PostV1Posts200Response{}, nil
}

// (DELETE /v1/posts/{post_id})
func (a *adapter) DeleteV1PostsPostId(
	ctx context.Context,
	request codegen.DeleteV1PostsPostIdRequestObject,
) (codegen.DeleteV1PostsPostIdResponseObject, error) {
	// TODO: use refresh token too
	// TODO: make it helper function
	_, userID, err := a.authService.ValidateAndRefresh(
		ctx,
		&model.TokenPair{AccessToken: request.Params.XSESSION, RefreshToken: ""},
	)
	switch {
	case errors.Is(err, service.ErrUnsupportedClaims) || errors.Is(err, service.ErrUnauthorized):
		return codegen.DeleteV1PostsPostId401Response{}, nil
	case err != nil:
		return nil, err
	}

	err = a.postsClient.DeletePost(ctx, request.PostId, *userID)
	if errors.Is(err, clients.ErrPostNotFound) {
		return codegen.DeleteV1PostsPostId404Response{}, nil
	}
	if errors.Is(err, clients.ErrInsufficientPermissions) {
		return codegen.DeleteV1PostsPostId403Response{}, nil
	}
	return codegen.DeleteV1PostsPostId200Response{}, nil
}

// (PATCH /v1/posts/{post_id})
func (a *adapter) PatchV1PostsPostId(
	ctx context.Context,
	request codegen.PatchV1PostsPostIdRequestObject,
) (codegen.PatchV1PostsPostIdResponseObject, error) {
	// TODO: use refresh token too
	// TODO: make it helper function
	_, userID, err := a.authService.ValidateAndRefresh(
		ctx,
		&model.TokenPair{AccessToken: request.Params.XSESSION, RefreshToken: ""},
	)
	switch {
	case errors.Is(err, service.ErrUnsupportedClaims) || errors.Is(err, service.ErrUnauthorized):
		return codegen.PatchV1PostsPostId401Response{}, nil
	case err != nil:
		return nil, err
	}

	err = a.postsClient.EditPost(ctx, request.PostId, *userID, request.Body.Content)
	if errors.Is(err, clients.ErrPostNotFound) {
		return codegen.PatchV1PostsPostId404Response{}, nil
	}
	if errors.Is(err, clients.ErrContentIsEmpty) {
		return codegen.PatchV1PostsPostId422JSONResponse(codegen.ErrorMessage{Error: err.Error()}), nil
	}
	if errors.Is(err, clients.ErrInsufficientPermissions) {
		return codegen.PatchV1PostsPostId403Response{}, nil
	}
	return codegen.PatchV1PostsPostId200Response{}, nil
}

// (GET /v1/posts/{post_id})
func (a *adapter) GetV1PostsPostId(
	ctx context.Context,
	request codegen.GetV1PostsPostIdRequestObject,
) (codegen.GetV1PostsPostIdResponseObject, error) {
	// TODO: use refresh token too
	// TODO: make it helper function
	_, _, err := a.authService.ValidateAndRefresh(
		ctx,
		&model.TokenPair{AccessToken: request.Params.XSESSION, RefreshToken: ""},
	)
	switch {
	case errors.Is(err, service.ErrUnsupportedClaims) || errors.Is(err, service.ErrUnauthorized):
		return codegen.GetV1PostsPostId401Response{}, nil
	case err != nil:
		return nil, err
	}

	post, err := a.postsClient.GetPost(ctx, request.PostId)
	if errors.Is(err, clients.ErrPostNotFound) {
		return codegen.GetV1PostsPostId404Response{}, nil
	}
	return codegen.GetV1PostsPostId200JSONResponse(
		codegen.Post{Id: post.ID, AuthorId: post.AuthorID.String(), Content: post.Content},
	), nil
}

// (POST /v1/posts/list)
func (a *adapter) PostV1PostsList(
	ctx context.Context,
	request codegen.PostV1PostsListRequestObject,
) (codegen.PostV1PostsListResponseObject, error) {
	// TODO: use refresh token too
	// TODO: make it helper function
	_, _, err := a.authService.ValidateAndRefresh(
		ctx,
		&model.TokenPair{AccessToken: request.Params.XSESSION, RefreshToken: ""},
	)
	switch {
	case errors.Is(err, service.ErrUnsupportedClaims) || errors.Is(err, service.ErrUnauthorized):
		return codegen.PostV1PostsList401Response{}, nil
	case err != nil:
		return nil, err
	}

	var pageToken string
	if request.Params.PageToken != nil {
		pageToken = *request.Params.PageToken
	} else {
		pageToken = ""
	}

	posts, nextPageToken, err := a.postsClient.ListPosts(
		ctx,
		uint32(request.Params.PageSize),
		pageToken,
	)
	if errors.Is(err, clients.ErrBadPageToken) {
		return codegen.PostV1PostsList422JSONResponse(codegen.ErrorMessage{Error: err.Error()}), nil
	}
	if err != nil {
		return nil, err
	}

	respPosts := make([]codegen.Post, 0, len(posts))
	for _, post := range posts {
		respPosts = append(
			respPosts,
			codegen.Post{Id: post.ID, AuthorId: post.AuthorID.String(), Content: post.Content},
		)
	}
	return codegen.PostV1PostsList200JSONResponse{NextPageToken: nextPageToken, Posts: respPosts}, nil
}

func (a *adapter) Serve() error {
	logger, err := logger.GetLogger(true)
	if err != nil {
		return err
	}
	a.logger = logger

	handlerOpts := codegen.StrictHTTPServerOptions{
		RequestErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		},
		ResponseErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
			logger.Error(err.Error())
			http.Error(w, "server got itself in trouble", http.StatusInternalServerError)
		},
	}
	strictHandler := codegen.NewStrictHandlerWithOptions(a, make([]codegen.StrictMiddlewareFunc, 0), handlerOpts)
	options := codegen.ChiServerOptions{
		BaseURL:     a.cfg.BasePath,
		BaseRouter:  chi.NewRouter(),
		Middlewares: make([]codegen.MiddlewareFunc, 0),
	}

	a.server = &http.Server{
		Handler: codegen.HandlerWithOptions(strictHandler, options),
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
	userService service.User,
	postsClient clients.PostsClient) Adapter {
	return &adapter{
		cfg:         config,
		authService: authService,
		userService: userService,
		postsClient: postsClient,
	}
}
