// Package codegen provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen/v2 version v2.0.0 DO NOT EDIT.
package codegen

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/oapi-codegen/runtime"
	strictnethttp "github.com/oapi-codegen/runtime/strictmiddleware/nethttp"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

// ErrorMessage defines model for ErrorMessage.
type ErrorMessage struct {
	Error string `json:"error"`
}

// Post defines model for Post.
type Post struct {
	AuthorId openapi_types.UUID `json:"author_id"`
	Content  string             `json:"content"`
	Id       openapi_types.UUID `json:"id"`
}

// PostV1AuthJSONBody defines parameters for PostV1Auth.
type PostV1AuthJSONBody struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// PostV1PostsJSONBody defines parameters for PostV1Posts.
type PostV1PostsJSONBody struct {
	Content string `json:"content"`
}

// PostV1PostsParams defines parameters for PostV1Posts.
type PostV1PostsParams struct {
	XSESSION string `form:"X_SESSION" json:"X_SESSION"`
}

// PostV1PostsListParams defines parameters for PostV1PostsList.
type PostV1PostsListParams struct {
	PageSize  int     `form:"page_size" json:"page_size"`
	PageToken *string `form:"page_token,omitempty" json:"page_token,omitempty"`
	XSESSION  string  `form:"X_SESSION" json:"X_SESSION"`
}

// DeleteV1PostsPostIdParams defines parameters for DeleteV1PostsPostId.
type DeleteV1PostsPostIdParams struct {
	XSESSION string `form:"X_SESSION" json:"X_SESSION"`
}

// GetV1PostsPostIdParams defines parameters for GetV1PostsPostId.
type GetV1PostsPostIdParams struct {
	XSESSION string `form:"X_SESSION" json:"X_SESSION"`
}

// PatchV1PostsPostIdJSONBody defines parameters for PatchV1PostsPostId.
type PatchV1PostsPostIdJSONBody struct {
	Content string `json:"content"`
}

// PatchV1PostsPostIdParams defines parameters for PatchV1PostsPostId.
type PatchV1PostsPostIdParams struct {
	XSESSION string `form:"X_SESSION" json:"X_SESSION"`
}

// PostV1UsersJSONBody defines parameters for PostV1Users.
type PostV1UsersJSONBody struct {
	Email    *string `json:"email,omitempty"`
	Login    string  `json:"login"`
	Name     *string `json:"name,omitempty"`
	Password string  `json:"password"`
	Phone    *string `json:"phone,omitempty"`
	Surname  *string `json:"surname,omitempty"`
}

// PatchV1UsersUserIdJSONBody defines parameters for PatchV1UsersUserId.
type PatchV1UsersUserIdJSONBody struct {
	Email   *string `json:"email,omitempty"`
	Name    *string `json:"name,omitempty"`
	Phone   *string `json:"phone,omitempty"`
	Surname *string `json:"surname,omitempty"`
}

// PatchV1UsersUserIdParams defines parameters for PatchV1UsersUserId.
type PatchV1UsersUserIdParams struct {
	XSESSION string `form:"X_SESSION" json:"X_SESSION"`
}

// PostV1AuthJSONRequestBody defines body for PostV1Auth for application/json ContentType.
type PostV1AuthJSONRequestBody PostV1AuthJSONBody

// PostV1PostsJSONRequestBody defines body for PostV1Posts for application/json ContentType.
type PostV1PostsJSONRequestBody PostV1PostsJSONBody

// PatchV1PostsPostIdJSONRequestBody defines body for PatchV1PostsPostId for application/json ContentType.
type PatchV1PostsPostIdJSONRequestBody PatchV1PostsPostIdJSONBody

// PostV1UsersJSONRequestBody defines body for PostV1Users for application/json ContentType.
type PostV1UsersJSONRequestBody PostV1UsersJSONBody

// PatchV1UsersUserIdJSONRequestBody defines body for PatchV1UsersUserId for application/json ContentType.
type PatchV1UsersUserIdJSONRequestBody PatchV1UsersUserIdJSONBody

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Authentificate user
	// (POST /v1/auth)
	PostV1Auth(w http.ResponseWriter, r *http.Request)

	// (POST /v1/posts)
	PostV1Posts(w http.ResponseWriter, r *http.Request, params PostV1PostsParams)

	// (POST /v1/posts/list)
	PostV1PostsList(w http.ResponseWriter, r *http.Request, params PostV1PostsListParams)

	// (DELETE /v1/posts/{post_id})
	DeleteV1PostsPostId(w http.ResponseWriter, r *http.Request, postId openapi_types.UUID, params DeleteV1PostsPostIdParams)

	// (GET /v1/posts/{post_id})
	GetV1PostsPostId(w http.ResponseWriter, r *http.Request, postId openapi_types.UUID, params GetV1PostsPostIdParams)

	// (PATCH /v1/posts/{post_id})
	PatchV1PostsPostId(w http.ResponseWriter, r *http.Request, postId openapi_types.UUID, params PatchV1PostsPostIdParams)

	// (POST /v1/users)
	PostV1Users(w http.ResponseWriter, r *http.Request)

	// (PATCH /v1/users/{user_id})
	PatchV1UsersUserId(w http.ResponseWriter, r *http.Request, userId openapi_types.UUID, params PatchV1UsersUserIdParams)
}

// Unimplemented server implementation that returns http.StatusNotImplemented for each endpoint.

type Unimplemented struct{}

// Authentificate user
// (POST /v1/auth)
func (_ Unimplemented) PostV1Auth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

// (POST /v1/posts)
func (_ Unimplemented) PostV1Posts(w http.ResponseWriter, r *http.Request, params PostV1PostsParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

// (POST /v1/posts/list)
func (_ Unimplemented) PostV1PostsList(w http.ResponseWriter, r *http.Request, params PostV1PostsListParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

// (DELETE /v1/posts/{post_id})
func (_ Unimplemented) DeleteV1PostsPostId(w http.ResponseWriter, r *http.Request, postId openapi_types.UUID, params DeleteV1PostsPostIdParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

// (GET /v1/posts/{post_id})
func (_ Unimplemented) GetV1PostsPostId(w http.ResponseWriter, r *http.Request, postId openapi_types.UUID, params GetV1PostsPostIdParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

// (PATCH /v1/posts/{post_id})
func (_ Unimplemented) PatchV1PostsPostId(w http.ResponseWriter, r *http.Request, postId openapi_types.UUID, params PatchV1PostsPostIdParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

// (POST /v1/users)
func (_ Unimplemented) PostV1Users(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

// (PATCH /v1/users/{user_id})
func (_ Unimplemented) PatchV1UsersUserId(w http.ResponseWriter, r *http.Request, userId openapi_types.UUID, params PatchV1UsersUserIdParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []MiddlewareFunc
	ErrorHandlerFunc   func(w http.ResponseWriter, r *http.Request, err error)
}

type MiddlewareFunc func(http.Handler) http.Handler

// PostV1Auth operation middleware
func (siw *ServerInterfaceWrapper) PostV1Auth(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.PostV1Auth(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// PostV1Posts operation middleware
func (siw *ServerInterfaceWrapper) PostV1Posts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// Parameter object where we will unmarshal all parameters from the context
	var params PostV1PostsParams

	var cookie *http.Cookie

	if cookie, err = r.Cookie("X_SESSION"); err == nil {
		var value string
		err = runtime.BindStyledParameter("simple", true, "X_SESSION", cookie.Value, &value)
		if err != nil {
			siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "X_SESSION", Err: err})
			return
		}
		params.XSESSION = value

	} else {
		siw.ErrorHandlerFunc(w, r, &RequiredParamError{ParamName: "X_SESSION"})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.PostV1Posts(w, r, params)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// PostV1PostsList operation middleware
func (siw *ServerInterfaceWrapper) PostV1PostsList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// Parameter object where we will unmarshal all parameters from the context
	var params PostV1PostsListParams

	// ------------- Required query parameter "page_size" -------------

	if paramValue := r.URL.Query().Get("page_size"); paramValue != "" {

	} else {
		siw.ErrorHandlerFunc(w, r, &RequiredParamError{ParamName: "page_size"})
		return
	}

	err = runtime.BindQueryParameter("form", true, true, "page_size", r.URL.Query(), &params.PageSize)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "page_size", Err: err})
		return
	}

	// ------------- Optional query parameter "page_token" -------------

	err = runtime.BindQueryParameter("form", true, false, "page_token", r.URL.Query(), &params.PageToken)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "page_token", Err: err})
		return
	}

	var cookie *http.Cookie

	if cookie, err = r.Cookie("X_SESSION"); err == nil {
		var value string
		err = runtime.BindStyledParameter("simple", true, "X_SESSION", cookie.Value, &value)
		if err != nil {
			siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "X_SESSION", Err: err})
			return
		}
		params.XSESSION = value

	} else {
		siw.ErrorHandlerFunc(w, r, &RequiredParamError{ParamName: "X_SESSION"})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.PostV1PostsList(w, r, params)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// DeleteV1PostsPostId operation middleware
func (siw *ServerInterfaceWrapper) DeleteV1PostsPostId(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "post_id" -------------
	var postId openapi_types.UUID

	err = runtime.BindStyledParameterWithLocation("simple", false, "post_id", runtime.ParamLocationPath, chi.URLParam(r, "post_id"), &postId)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "post_id", Err: err})
		return
	}

	// Parameter object where we will unmarshal all parameters from the context
	var params DeleteV1PostsPostIdParams

	var cookie *http.Cookie

	if cookie, err = r.Cookie("X_SESSION"); err == nil {
		var value string
		err = runtime.BindStyledParameter("simple", true, "X_SESSION", cookie.Value, &value)
		if err != nil {
			siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "X_SESSION", Err: err})
			return
		}
		params.XSESSION = value

	} else {
		siw.ErrorHandlerFunc(w, r, &RequiredParamError{ParamName: "X_SESSION"})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.DeleteV1PostsPostId(w, r, postId, params)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// GetV1PostsPostId operation middleware
func (siw *ServerInterfaceWrapper) GetV1PostsPostId(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "post_id" -------------
	var postId openapi_types.UUID

	err = runtime.BindStyledParameterWithLocation("simple", false, "post_id", runtime.ParamLocationPath, chi.URLParam(r, "post_id"), &postId)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "post_id", Err: err})
		return
	}

	// Parameter object where we will unmarshal all parameters from the context
	var params GetV1PostsPostIdParams

	var cookie *http.Cookie

	if cookie, err = r.Cookie("X_SESSION"); err == nil {
		var value string
		err = runtime.BindStyledParameter("simple", true, "X_SESSION", cookie.Value, &value)
		if err != nil {
			siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "X_SESSION", Err: err})
			return
		}
		params.XSESSION = value

	} else {
		siw.ErrorHandlerFunc(w, r, &RequiredParamError{ParamName: "X_SESSION"})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetV1PostsPostId(w, r, postId, params)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// PatchV1PostsPostId operation middleware
func (siw *ServerInterfaceWrapper) PatchV1PostsPostId(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "post_id" -------------
	var postId openapi_types.UUID

	err = runtime.BindStyledParameterWithLocation("simple", false, "post_id", runtime.ParamLocationPath, chi.URLParam(r, "post_id"), &postId)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "post_id", Err: err})
		return
	}

	// Parameter object where we will unmarshal all parameters from the context
	var params PatchV1PostsPostIdParams

	var cookie *http.Cookie

	if cookie, err = r.Cookie("X_SESSION"); err == nil {
		var value string
		err = runtime.BindStyledParameter("simple", true, "X_SESSION", cookie.Value, &value)
		if err != nil {
			siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "X_SESSION", Err: err})
			return
		}
		params.XSESSION = value

	} else {
		siw.ErrorHandlerFunc(w, r, &RequiredParamError{ParamName: "X_SESSION"})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.PatchV1PostsPostId(w, r, postId, params)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// PostV1Users operation middleware
func (siw *ServerInterfaceWrapper) PostV1Users(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.PostV1Users(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// PatchV1UsersUserId operation middleware
func (siw *ServerInterfaceWrapper) PatchV1UsersUserId(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "user_id" -------------
	var userId openapi_types.UUID

	err = runtime.BindStyledParameterWithLocation("simple", false, "user_id", runtime.ParamLocationPath, chi.URLParam(r, "user_id"), &userId)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "user_id", Err: err})
		return
	}

	// Parameter object where we will unmarshal all parameters from the context
	var params PatchV1UsersUserIdParams

	var cookie *http.Cookie

	if cookie, err = r.Cookie("X_SESSION"); err == nil {
		var value string
		err = runtime.BindStyledParameter("simple", true, "X_SESSION", cookie.Value, &value)
		if err != nil {
			siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "X_SESSION", Err: err})
			return
		}
		params.XSESSION = value

	} else {
		siw.ErrorHandlerFunc(w, r, &RequiredParamError{ParamName: "X_SESSION"})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.PatchV1UsersUserId(w, r, userId, params)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

type UnescapedCookieParamError struct {
	ParamName string
	Err       error
}

func (e *UnescapedCookieParamError) Error() string {
	return fmt.Sprintf("error unescaping cookie parameter '%s'", e.ParamName)
}

func (e *UnescapedCookieParamError) Unwrap() error {
	return e.Err
}

type UnmarshalingParamError struct {
	ParamName string
	Err       error
}

func (e *UnmarshalingParamError) Error() string {
	return fmt.Sprintf("Error unmarshaling parameter %s as JSON: %s", e.ParamName, e.Err.Error())
}

func (e *UnmarshalingParamError) Unwrap() error {
	return e.Err
}

type RequiredParamError struct {
	ParamName string
}

func (e *RequiredParamError) Error() string {
	return fmt.Sprintf("Query argument %s is required, but not found", e.ParamName)
}

type RequiredHeaderError struct {
	ParamName string
	Err       error
}

func (e *RequiredHeaderError) Error() string {
	return fmt.Sprintf("Header parameter %s is required, but not found", e.ParamName)
}

func (e *RequiredHeaderError) Unwrap() error {
	return e.Err
}

type InvalidParamFormatError struct {
	ParamName string
	Err       error
}

func (e *InvalidParamFormatError) Error() string {
	return fmt.Sprintf("Invalid format for parameter %s: %s", e.ParamName, e.Err.Error())
}

func (e *InvalidParamFormatError) Unwrap() error {
	return e.Err
}

type TooManyValuesForParamError struct {
	ParamName string
	Count     int
}

func (e *TooManyValuesForParamError) Error() string {
	return fmt.Sprintf("Expected one value for %s, got %d", e.ParamName, e.Count)
}

// Handler creates http.Handler with routing matching OpenAPI spec.
func Handler(si ServerInterface) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{})
}

type ChiServerOptions struct {
	BaseURL          string
	BaseRouter       chi.Router
	Middlewares      []MiddlewareFunc
	ErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)
}

// HandlerFromMux creates http.Handler with routing matching OpenAPI spec based on the provided mux.
func HandlerFromMux(si ServerInterface, r chi.Router) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{
		BaseRouter: r,
	})
}

func HandlerFromMuxWithBaseURL(si ServerInterface, r chi.Router, baseURL string) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{
		BaseURL:    baseURL,
		BaseRouter: r,
	})
}

// HandlerWithOptions creates http.Handler with additional options
func HandlerWithOptions(si ServerInterface, options ChiServerOptions) http.Handler {
	r := options.BaseRouter

	if r == nil {
		r = chi.NewRouter()
	}
	if options.ErrorHandlerFunc == nil {
		options.ErrorHandlerFunc = func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}
	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: options.Middlewares,
		ErrorHandlerFunc:   options.ErrorHandlerFunc,
	}

	r.Group(func(r chi.Router) {
		r.Post(options.BaseURL+"/v1/auth", wrapper.PostV1Auth)
	})
	r.Group(func(r chi.Router) {
		r.Post(options.BaseURL+"/v1/posts", wrapper.PostV1Posts)
	})
	r.Group(func(r chi.Router) {
		r.Post(options.BaseURL+"/v1/posts/list", wrapper.PostV1PostsList)
	})
	r.Group(func(r chi.Router) {
		r.Delete(options.BaseURL+"/v1/posts/{post_id}", wrapper.DeleteV1PostsPostId)
	})
	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/v1/posts/{post_id}", wrapper.GetV1PostsPostId)
	})
	r.Group(func(r chi.Router) {
		r.Patch(options.BaseURL+"/v1/posts/{post_id}", wrapper.PatchV1PostsPostId)
	})
	r.Group(func(r chi.Router) {
		r.Post(options.BaseURL+"/v1/users", wrapper.PostV1Users)
	})
	r.Group(func(r chi.Router) {
		r.Patch(options.BaseURL+"/v1/users/{user_id}", wrapper.PatchV1UsersUserId)
	})

	return r
}

type PostV1AuthRequestObject struct {
	Body *PostV1AuthJSONRequestBody
}

type PostV1AuthResponseObject interface {
	VisitPostV1AuthResponse(w http.ResponseWriter) error
}

type PostV1Auth200ResponseHeaders struct {
	SetCookie string
}

type PostV1Auth200Response struct {
	Headers PostV1Auth200ResponseHeaders
}

func (response PostV1Auth200Response) VisitPostV1AuthResponse(w http.ResponseWriter) error {
	w.Header().Set("Set-Cookie", fmt.Sprint(response.Headers.SetCookie))
	w.WriteHeader(200)
	return nil
}

type PostV1Auth401JSONResponse ErrorMessage

func (response PostV1Auth401JSONResponse) VisitPostV1AuthResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(401)

	return json.NewEncoder(w).Encode(response)
}

type PostV1PostsRequestObject struct {
	Params PostV1PostsParams
	Body   *PostV1PostsJSONRequestBody
}

type PostV1PostsResponseObject interface {
	VisitPostV1PostsResponse(w http.ResponseWriter) error
}

type PostV1Posts200Response struct {
}

func (response PostV1Posts200Response) VisitPostV1PostsResponse(w http.ResponseWriter) error {
	w.WriteHeader(200)
	return nil
}

type PostV1Posts401Response struct {
}

func (response PostV1Posts401Response) VisitPostV1PostsResponse(w http.ResponseWriter) error {
	w.WriteHeader(401)
	return nil
}

type PostV1Posts422JSONResponse ErrorMessage

func (response PostV1Posts422JSONResponse) VisitPostV1PostsResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(422)

	return json.NewEncoder(w).Encode(response)
}

type PostV1PostsListRequestObject struct {
	Params PostV1PostsListParams
}

type PostV1PostsListResponseObject interface {
	VisitPostV1PostsListResponse(w http.ResponseWriter) error
}

type PostV1PostsList200JSONResponse struct {
	NextPageToken string `json:"next_page_token"`
	Posts         []Post `json:"posts"`
}

func (response PostV1PostsList200JSONResponse) VisitPostV1PostsListResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response)
}

type PostV1PostsList401Response struct {
}

func (response PostV1PostsList401Response) VisitPostV1PostsListResponse(w http.ResponseWriter) error {
	w.WriteHeader(401)
	return nil
}

type PostV1PostsList422JSONResponse ErrorMessage

func (response PostV1PostsList422JSONResponse) VisitPostV1PostsListResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(422)

	return json.NewEncoder(w).Encode(response)
}

type DeleteV1PostsPostIdRequestObject struct {
	PostId openapi_types.UUID `json:"post_id"`
	Params DeleteV1PostsPostIdParams
}

type DeleteV1PostsPostIdResponseObject interface {
	VisitDeleteV1PostsPostIdResponse(w http.ResponseWriter) error
}

type DeleteV1PostsPostId200Response struct {
}

func (response DeleteV1PostsPostId200Response) VisitDeleteV1PostsPostIdResponse(w http.ResponseWriter) error {
	w.WriteHeader(200)
	return nil
}

type DeleteV1PostsPostId401Response struct {
}

func (response DeleteV1PostsPostId401Response) VisitDeleteV1PostsPostIdResponse(w http.ResponseWriter) error {
	w.WriteHeader(401)
	return nil
}

type DeleteV1PostsPostId403Response struct {
}

func (response DeleteV1PostsPostId403Response) VisitDeleteV1PostsPostIdResponse(w http.ResponseWriter) error {
	w.WriteHeader(403)
	return nil
}

type DeleteV1PostsPostId404Response struct {
}

func (response DeleteV1PostsPostId404Response) VisitDeleteV1PostsPostIdResponse(w http.ResponseWriter) error {
	w.WriteHeader(404)
	return nil
}

type GetV1PostsPostIdRequestObject struct {
	PostId openapi_types.UUID `json:"post_id"`
	Params GetV1PostsPostIdParams
}

type GetV1PostsPostIdResponseObject interface {
	VisitGetV1PostsPostIdResponse(w http.ResponseWriter) error
}

type GetV1PostsPostId200JSONResponse Post

func (response GetV1PostsPostId200JSONResponse) VisitGetV1PostsPostIdResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response)
}

type GetV1PostsPostId401Response struct {
}

func (response GetV1PostsPostId401Response) VisitGetV1PostsPostIdResponse(w http.ResponseWriter) error {
	w.WriteHeader(401)
	return nil
}

type GetV1PostsPostId404Response struct {
}

func (response GetV1PostsPostId404Response) VisitGetV1PostsPostIdResponse(w http.ResponseWriter) error {
	w.WriteHeader(404)
	return nil
}

type GetV1PostsPostId422JSONResponse ErrorMessage

func (response GetV1PostsPostId422JSONResponse) VisitGetV1PostsPostIdResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(422)

	return json.NewEncoder(w).Encode(response)
}

type PatchV1PostsPostIdRequestObject struct {
	PostId openapi_types.UUID `json:"post_id"`
	Params PatchV1PostsPostIdParams
	Body   *PatchV1PostsPostIdJSONRequestBody
}

type PatchV1PostsPostIdResponseObject interface {
	VisitPatchV1PostsPostIdResponse(w http.ResponseWriter) error
}

type PatchV1PostsPostId200Response struct {
}

func (response PatchV1PostsPostId200Response) VisitPatchV1PostsPostIdResponse(w http.ResponseWriter) error {
	w.WriteHeader(200)
	return nil
}

type PatchV1PostsPostId401Response struct {
}

func (response PatchV1PostsPostId401Response) VisitPatchV1PostsPostIdResponse(w http.ResponseWriter) error {
	w.WriteHeader(401)
	return nil
}

type PatchV1PostsPostId403Response struct {
}

func (response PatchV1PostsPostId403Response) VisitPatchV1PostsPostIdResponse(w http.ResponseWriter) error {
	w.WriteHeader(403)
	return nil
}

type PatchV1PostsPostId404Response struct {
}

func (response PatchV1PostsPostId404Response) VisitPatchV1PostsPostIdResponse(w http.ResponseWriter) error {
	w.WriteHeader(404)
	return nil
}

type PatchV1PostsPostId422JSONResponse ErrorMessage

func (response PatchV1PostsPostId422JSONResponse) VisitPatchV1PostsPostIdResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(422)

	return json.NewEncoder(w).Encode(response)
}

type PostV1UsersRequestObject struct {
	Body *PostV1UsersJSONRequestBody
}

type PostV1UsersResponseObject interface {
	VisitPostV1UsersResponse(w http.ResponseWriter) error
}

type PostV1Users200Response struct {
}

func (response PostV1Users200Response) VisitPostV1UsersResponse(w http.ResponseWriter) error {
	w.WriteHeader(200)
	return nil
}

type PostV1Users422JSONResponse ErrorMessage

func (response PostV1Users422JSONResponse) VisitPostV1UsersResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(422)

	return json.NewEncoder(w).Encode(response)
}

type PatchV1UsersUserIdRequestObject struct {
	UserId openapi_types.UUID `json:"user_id"`
	Params PatchV1UsersUserIdParams
	Body   *PatchV1UsersUserIdJSONRequestBody
}

type PatchV1UsersUserIdResponseObject interface {
	VisitPatchV1UsersUserIdResponse(w http.ResponseWriter) error
}

type PatchV1UsersUserId200Response struct {
}

func (response PatchV1UsersUserId200Response) VisitPatchV1UsersUserIdResponse(w http.ResponseWriter) error {
	w.WriteHeader(200)
	return nil
}

type PatchV1UsersUserId401Response struct {
}

func (response PatchV1UsersUserId401Response) VisitPatchV1UsersUserIdResponse(w http.ResponseWriter) error {
	w.WriteHeader(401)
	return nil
}

type PatchV1UsersUserId422JSONResponse ErrorMessage

func (response PatchV1UsersUserId422JSONResponse) VisitPatchV1UsersUserIdResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(422)

	return json.NewEncoder(w).Encode(response)
}

// StrictServerInterface represents all server handlers.
type StrictServerInterface interface {
	// Authentificate user
	// (POST /v1/auth)
	PostV1Auth(ctx context.Context, request PostV1AuthRequestObject) (PostV1AuthResponseObject, error)

	// (POST /v1/posts)
	PostV1Posts(ctx context.Context, request PostV1PostsRequestObject) (PostV1PostsResponseObject, error)

	// (POST /v1/posts/list)
	PostV1PostsList(ctx context.Context, request PostV1PostsListRequestObject) (PostV1PostsListResponseObject, error)

	// (DELETE /v1/posts/{post_id})
	DeleteV1PostsPostId(ctx context.Context, request DeleteV1PostsPostIdRequestObject) (DeleteV1PostsPostIdResponseObject, error)

	// (GET /v1/posts/{post_id})
	GetV1PostsPostId(ctx context.Context, request GetV1PostsPostIdRequestObject) (GetV1PostsPostIdResponseObject, error)

	// (PATCH /v1/posts/{post_id})
	PatchV1PostsPostId(ctx context.Context, request PatchV1PostsPostIdRequestObject) (PatchV1PostsPostIdResponseObject, error)

	// (POST /v1/users)
	PostV1Users(ctx context.Context, request PostV1UsersRequestObject) (PostV1UsersResponseObject, error)

	// (PATCH /v1/users/{user_id})
	PatchV1UsersUserId(ctx context.Context, request PatchV1UsersUserIdRequestObject) (PatchV1UsersUserIdResponseObject, error)
}

type StrictHandlerFunc = strictnethttp.StrictHttpHandlerFunc
type StrictMiddlewareFunc = strictnethttp.StrictHttpMiddlewareFunc

type StrictHTTPServerOptions struct {
	RequestErrorHandlerFunc  func(w http.ResponseWriter, r *http.Request, err error)
	ResponseErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)
}

func NewStrictHandler(ssi StrictServerInterface, middlewares []StrictMiddlewareFunc) ServerInterface {
	return &strictHandler{ssi: ssi, middlewares: middlewares, options: StrictHTTPServerOptions{
		RequestErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		},
		ResponseErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		},
	}}
}

func NewStrictHandlerWithOptions(ssi StrictServerInterface, middlewares []StrictMiddlewareFunc, options StrictHTTPServerOptions) ServerInterface {
	return &strictHandler{ssi: ssi, middlewares: middlewares, options: options}
}

type strictHandler struct {
	ssi         StrictServerInterface
	middlewares []StrictMiddlewareFunc
	options     StrictHTTPServerOptions
}

// PostV1Auth operation middleware
func (sh *strictHandler) PostV1Auth(w http.ResponseWriter, r *http.Request) {
	var request PostV1AuthRequestObject

	var body PostV1AuthJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		sh.options.RequestErrorHandlerFunc(w, r, fmt.Errorf("can't decode JSON body: %w", err))
		return
	}
	request.Body = &body

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
		return sh.ssi.PostV1Auth(ctx, request.(PostV1AuthRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "PostV1Auth")
	}

	response, err := handler(r.Context(), w, r, request)

	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(PostV1AuthResponseObject); ok {
		if err := validResponse.VisitPostV1AuthResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("unexpected response type: %T", response))
	}
}

// PostV1Posts operation middleware
func (sh *strictHandler) PostV1Posts(w http.ResponseWriter, r *http.Request, params PostV1PostsParams) {
	var request PostV1PostsRequestObject

	request.Params = params

	var body PostV1PostsJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		sh.options.RequestErrorHandlerFunc(w, r, fmt.Errorf("can't decode JSON body: %w", err))
		return
	}
	request.Body = &body

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
		return sh.ssi.PostV1Posts(ctx, request.(PostV1PostsRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "PostV1Posts")
	}

	response, err := handler(r.Context(), w, r, request)

	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(PostV1PostsResponseObject); ok {
		if err := validResponse.VisitPostV1PostsResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("unexpected response type: %T", response))
	}
}

// PostV1PostsList operation middleware
func (sh *strictHandler) PostV1PostsList(w http.ResponseWriter, r *http.Request, params PostV1PostsListParams) {
	var request PostV1PostsListRequestObject

	request.Params = params

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
		return sh.ssi.PostV1PostsList(ctx, request.(PostV1PostsListRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "PostV1PostsList")
	}

	response, err := handler(r.Context(), w, r, request)

	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(PostV1PostsListResponseObject); ok {
		if err := validResponse.VisitPostV1PostsListResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("unexpected response type: %T", response))
	}
}

// DeleteV1PostsPostId operation middleware
func (sh *strictHandler) DeleteV1PostsPostId(w http.ResponseWriter, r *http.Request, postId openapi_types.UUID, params DeleteV1PostsPostIdParams) {
	var request DeleteV1PostsPostIdRequestObject

	request.PostId = postId
	request.Params = params

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
		return sh.ssi.DeleteV1PostsPostId(ctx, request.(DeleteV1PostsPostIdRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "DeleteV1PostsPostId")
	}

	response, err := handler(r.Context(), w, r, request)

	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(DeleteV1PostsPostIdResponseObject); ok {
		if err := validResponse.VisitDeleteV1PostsPostIdResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("unexpected response type: %T", response))
	}
}

// GetV1PostsPostId operation middleware
func (sh *strictHandler) GetV1PostsPostId(w http.ResponseWriter, r *http.Request, postId openapi_types.UUID, params GetV1PostsPostIdParams) {
	var request GetV1PostsPostIdRequestObject

	request.PostId = postId
	request.Params = params

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
		return sh.ssi.GetV1PostsPostId(ctx, request.(GetV1PostsPostIdRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "GetV1PostsPostId")
	}

	response, err := handler(r.Context(), w, r, request)

	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(GetV1PostsPostIdResponseObject); ok {
		if err := validResponse.VisitGetV1PostsPostIdResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("unexpected response type: %T", response))
	}
}

// PatchV1PostsPostId operation middleware
func (sh *strictHandler) PatchV1PostsPostId(w http.ResponseWriter, r *http.Request, postId openapi_types.UUID, params PatchV1PostsPostIdParams) {
	var request PatchV1PostsPostIdRequestObject

	request.PostId = postId
	request.Params = params

	var body PatchV1PostsPostIdJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		sh.options.RequestErrorHandlerFunc(w, r, fmt.Errorf("can't decode JSON body: %w", err))
		return
	}
	request.Body = &body

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
		return sh.ssi.PatchV1PostsPostId(ctx, request.(PatchV1PostsPostIdRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "PatchV1PostsPostId")
	}

	response, err := handler(r.Context(), w, r, request)

	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(PatchV1PostsPostIdResponseObject); ok {
		if err := validResponse.VisitPatchV1PostsPostIdResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("unexpected response type: %T", response))
	}
}

// PostV1Users operation middleware
func (sh *strictHandler) PostV1Users(w http.ResponseWriter, r *http.Request) {
	var request PostV1UsersRequestObject

	var body PostV1UsersJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		sh.options.RequestErrorHandlerFunc(w, r, fmt.Errorf("can't decode JSON body: %w", err))
		return
	}
	request.Body = &body

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
		return sh.ssi.PostV1Users(ctx, request.(PostV1UsersRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "PostV1Users")
	}

	response, err := handler(r.Context(), w, r, request)

	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(PostV1UsersResponseObject); ok {
		if err := validResponse.VisitPostV1UsersResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("unexpected response type: %T", response))
	}
}

// PatchV1UsersUserId operation middleware
func (sh *strictHandler) PatchV1UsersUserId(w http.ResponseWriter, r *http.Request, userId openapi_types.UUID, params PatchV1UsersUserIdParams) {
	var request PatchV1UsersUserIdRequestObject

	request.UserId = userId
	request.Params = params

	var body PatchV1UsersUserIdJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		sh.options.RequestErrorHandlerFunc(w, r, fmt.Errorf("can't decode JSON body: %w", err))
		return
	}
	request.Body = &body

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
		return sh.ssi.PatchV1UsersUserId(ctx, request.(PatchV1UsersUserIdRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "PatchV1UsersUserId")
	}

	response, err := handler(r.Context(), w, r, request)

	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(PatchV1UsersUserIdResponseObject); ok {
		if err := validResponse.VisitPatchV1UsersUserIdResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("unexpected response type: %T", response))
	}
}
