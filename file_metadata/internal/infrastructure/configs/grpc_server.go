package configs

import (
	"context"
	"elex_storage/pkg/shared_kernel/models"

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
) {
	grpc_service.RegisterFileMetadataServiceServer(grpcServer, fileMetadataService)

	// Hook into the Fx Application Lifecycle
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			grpcPort := fmt.Sprintf(":%s", cfg.FileMetadataGrpcPort)
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
			fmt.Println("Shutting down gRPC server gracefully...")
			grpcServer.GracefulStop()
			return nil
		},
	})
}
