package grpcadapter

import (
	"context"
	"fmt"
	"net"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/aedobrynin/soa-hw/posts/internal/grpcadapter/postgrpc"
	"github.com/aedobrynin/soa-hw/posts/internal/service"
)

type Adapter struct {
	logger     *zap.Logger
	gRPCServer *grpc.Server
	config     *GRPCConfig
}

func New(
	logger *zap.Logger,
	post service.Post,
	config *GRPCConfig,
) *Adapter {
	loggingOpts := []logging.Option{
		logging.WithLogOnEvents(
			logging.StartCall, logging.FinishCall,
			logging.PayloadReceived, logging.PayloadSent,
		),
		// Add any other option (check functions starting with logging.With).
	}

	recoveryOpts := []recovery.Option{
		recovery.WithRecoveryHandler(func(p interface{}) (err error) {
			// TODO: log panic
			defer logger.Sync()
			logger.Error("Recovered from panic")
			return status.Errorf(codes.Internal, "internal error")
		}),
	}

	gRPCServer := grpc.NewServer(grpc.ChainUnaryInterceptor(
		recovery.UnaryServerInterceptor(recoveryOpts...),
		logging.UnaryServerInterceptor(InterceptorLogger(logger), loggingOpts...),
	))

	postgrpc.Register(gRPCServer, post)

	return &Adapter{
		logger:     logger,
		gRPCServer: gRPCServer,
		config:     config,
	}
}

func InterceptorLogger(l *zap.Logger) logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		var level zapcore.Level
		switch lvl {
		case logging.LevelError:
			level = zapcore.ErrorLevel
		case logging.LevelWarn:
			level = zapcore.WarnLevel
		case logging.LevelInfo:
			level = zapcore.InfoLevel
		case logging.LevelDebug:
			level = zapcore.DebugLevel
		}
		l.Log(level, msg)
		// TODO: fields
	})
}

func (a *Adapter) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

// Run runs gRPC server.
func (a *Adapter) Run() error {
	const op = "grpcapp.Run"

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.config.Port))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	a.logger.Sugar().Infof("gRPC server started at addr=%s", l.Addr().String())

	if err := a.gRPCServer.Serve(l); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// Stop stops gRPC server.
func (a *Adapter) Stop() {
	defer a.logger.Sync()
	a.logger.Info("stopping gRPC server")
	a.gRPCServer.GracefulStop()
}
