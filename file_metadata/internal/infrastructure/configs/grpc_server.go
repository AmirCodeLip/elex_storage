package configs

import (
	"context"
	"elex_storage/pkg/logger"
	"elex_storage/pkg/shared_kernel/models"
	"errors"
	"time"

	"elex_storage/pkg/shared_kernel"
	"elex_storage/pkg/shared_kernel/grpc_service"
	"fmt"
	"net"

	"go.uber.org/fx"
	"google.golang.org/grpc"
)

func NewGRPCServer() *grpc.Server {
	return grpc.NewServer()
}

func RegisterAndStartGRPCServer(
	lc fx.Lifecycle,
	grpcServer *grpc.Server,
	fileMetadataService grpc_service.FileMetadataServiceServer,
	cfg *models.ConfigEnv,
	l logger.Logger,
) error {
	grpcUrl, err := shared_kernel.ParseUrl(cfg.Server.GRPCListenUrl)
	if err != nil {
		return errors.New("GRPCListenUrl is not set or is invalid; set it in config.yml")
	}
	grpc_service.RegisterFileMetadataServiceServer(grpcServer, fileMetadataService)

	// Hook into the Fx Application Lifecycle
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			grpcPort := fmt.Sprintf(":%d", grpcUrl.Port)
			listener, err := net.Listen("tcp", grpcPort)
			if err != nil {
				return err
			}

			fmt.Printf("gRPC server listening at %v", listener.Addr())

			// grpcServer.Serve is blocking, so we run it in a goroutine
			go func() {
				if err := grpcServer.Serve(listener); err != nil {
					fmt.Printf("failed to serve: %v", err)
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
