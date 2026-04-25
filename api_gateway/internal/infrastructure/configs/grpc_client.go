package configs

import (
	"elex_storage/pkg/shared_kernel/grpc_service"
	"elex_storage/pkg/shared_kernel/models"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewGRPCClient(cfg *models.ConfigEnv) (grpc_service.UserServiceClient, grpc_service.FileMetadataServiceClient, error) {
	identityConn, err := grpc.NewClient(
		cfg.IdentityServiceGrpcAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, nil, err
	}
	fileMetadata, err := grpc.NewClient(
		cfg.FileMetadataGrpcAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, nil, err
	}

	userServiceClient := grpc_service.NewUserServiceClient(identityConn)
	fileMetadataClient := grpc_service.NewFileMetadataServiceClient(fileMetadata)
	return userServiceClient, fileMetadataClient, nil
}
