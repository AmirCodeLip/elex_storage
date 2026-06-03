package configs

import (
	"context"
	"elex_storage/pkg/logger"
	"elex_storage/pkg/shared_kernel"
	"elex_storage/pkg/shared_kernel/grpc_service"
	"elex_storage/pkg/shared_kernel/models"
	"errors"
	"fmt"
	"net"
	"time"

	"go.uber.org/fx"
	"google.golang.org/grpc"
)

func NewGRPCServer() *grpc.Server {
	return grpc.NewServer()
}

func RegisterAndStartGRPCServer(
	lc fx.Lifecycle,
	grpcServer *grpc.Server,
	userService grpc_service.UserServiceServer,
	cfg *models.ConfigEnv,
	l logger.Logger,
) error {
	grpcUrl, err := shared_kernel.ParseUrl(cfg.Server.GRPCListenUrl)
	if err != nil {
		return errors.New("GRPCListenUrl is not set or is invalid; set it in config.yml")
	}
	grpc_service.RegisterUserServiceServer(grpcServer, userService)

	// Hook into the Fx Application Lifecycle
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			listener, err := net.Listen("tcp", grpcUrl.Address)
			if err != nil {
				return err
			}
			fmt.Printf("GRPC server listening at %v", grpcUrl.FullAddress)
			l.Info(fmt.Sprintf("GRPC server listening at %s", grpcUrl.FullAddress))
			// grpcServer.Serve is blocking, so we run it in a goroutine
			go func() {
				if err := grpcServer.Serve(listener); err != nil {
					fmt.Printf("failed to serve: %v \n", err)
					l.Error(fmt.Sprintf("GRPC server have problem on listening %s", grpcUrl.FullAddress))
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			shutdownTime := time.Now()
			duration := time.Since(shutdownTime)
			fmt.Println("Shutting down gRPC server gracefully...")
			l.Info("Server shut down gracefully",
				"time", shutdownTime.Format("2006-01-02 15:04:05"),
				"duration", duration.Round(time.Millisecond),
			)
			grpcServer.GracefulStop()
			return nil
		},
	})
	return nil
}
