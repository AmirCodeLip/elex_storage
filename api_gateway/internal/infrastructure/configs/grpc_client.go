package configs

import (
	"elex_storage/pkg/shared_kernel/grpc_service"
	"elex_storage/pkg/shared_kernel/models"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewGRPCClient(cfg *models.ConfigEnv) (grpc_service.UserServiceClient, error) {
	identityConn, err := grpc.NewClient(
		cfg.IdentityServiceGrpcAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}
	client := grpc_service.NewUserServiceClient(identityConn)
	return client, nil
}
