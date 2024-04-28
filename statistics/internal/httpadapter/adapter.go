package httpadapter

import (
	"context"
	"net"
	"net/http"

	"github.com/aedobrynin/soa-hw/statistics/internal/httpadapter/codegen"
	"github.com/aedobrynin/soa-hw/statistics/internal/logger"

	"github.com/go-chi/chi/v5"
	"github.com/juju/zaputil/zapctx"
	"go.uber.org/zap"
)

type adapter struct {
	cfg *Config

	server *http.Server

	logger *zap.Logger
}

// (GET /v1/temp)
func (a *adapter) GetV1Temp(
	ctx context.Context,
	request codegen.GetV1TempRequestObject,
) (codegen.GetV1TempResponseObject, error) {
	return codegen.GetV1Temp200JSONResponse{}, nil
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
	config *Config) Adapter {
	return &adapter{
		cfg: config,
	}
}
