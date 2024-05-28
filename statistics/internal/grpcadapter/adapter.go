package grpcadapter

import (
	"fmt"
	"net"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/aedobrynin/soa-hw/statistics/internal/grpcadapter/statisticsgrpc"
	"github.com/aedobrynin/soa-hw/statistics/internal/service"
)

type Adapter struct {
	logger     *zap.Logger
	gRPCServer *grpc.Server
	config     *GRPCConfig
}

func New(
	logger *zap.Logger,
	statistics service.Statistics,
	config *GRPCConfig,
) *Adapter {
	recoveryOpts := []recovery.Option{
		recovery.WithRecoveryHandler(func(p interface{}) (err error) {
			// TODO: log panic
			defer func() {
				_ = logger.Sync()
			}()
			logger.Error("Recovered from panic")
			return status.Errorf(codes.Internal, "internal error")
		}),
	}

	gRPCServer := grpc.NewServer(grpc.ChainUnaryInterceptor(
		recovery.UnaryServerInterceptor(recoveryOpts...),
	))

	statisticsgrpc.Register(gRPCServer, statistics)

	return &Adapter{
		logger:     logger,
		gRPCServer: gRPCServer,
		config:     config,
	}
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
	defer func() {
		_ = a.logger.Sync()
	}()
	a.logger.Info("stopping gRPC server")
	a.gRPCServer.GracefulStop()
}
