package grpc_server

import (
	"context"
	"elex_storage/file_metadata/internal/use_case/cqrs/command_handlers"
	"elex_storage/file_metadata/internal/use_case/cqrs/commands"
	"elex_storage/pkg/logger"
	"elex_storage/pkg/shared_kernel"
	"elex_storage/pkg/shared_kernel/grpc_service"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var InvalidContentErr = errors.New("Invalid content is passed")

type FileMetadataService struct {
	grpc_service.UnimplementedFileMetadataServiceServer
	createDirectoryHandler *command_handlers.CreateDirectoryHandler
	getDirectoriesHandler  *command_handlers.GetDirectoriesHandler
	getDriveHandler        *command_handlers.GetDriveHandler
	logger                 logger.Logger
}

func NewFileMetadataService(getDirectoriesHandler *command_handlers.GetDirectoriesHandler, createDirectoryHandler *command_handlers.CreateDirectoryHandler, logger logger.Logger) grpc_service.FileMetadataServiceServer {
	return &FileMetadataService{getDirectoriesHandler: getDirectoriesHandler, createDirectoryHandler: createDirectoryHandler, logger: logger}
}

func (f *FileMetadataService) GetDirectories(context.Context, *grpc_service.GetDirectoriesRequest) (*grpc_service.StorageItemsResponse, error) {
	dirs, err := f.getDirectoriesHandler.Handle()
	if err != nil {
		f.logger.Error(err.Error())
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	response := grpc_service.StorageItemsResponse{}
	err = shared_kernel.MapToGrpc(&response.Directories, &dirs)
	if err != nil {
		f.logger.Error(err.Error())
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return &response, nil
}

func (f *FileMetadataService) CreateDirectory(ctx context.Context, req *grpc_service.CreateDirectoryRequest) (*grpc_service.DirectoryInfo, error) {
	var createDirectoryCommand commands.CreateDirectoryCommand
	err := shared_kernel.MapToGrpc(&createDirectoryCommand, &req)
	if err != nil {
		f.logger.Error(err.Error())
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	err, directory := f.createDirectoryHandler.Handle(createDirectoryCommand)
	if err != nil {
		return nil, status.Error(codes.Aborted, err.Error())
	}
	var response grpc_service.DirectoryInfo
	err = shared_kernel.MapFromGrpc(&response, &directory)
	if err != nil {
		f.logger.Error(err.Error())
		return nil, status.Error(codes.Aborted, err.Error())
	}
	return &response, nil
}

func (f *FileMetadataService) GetStorageItems(ctx context.Context, req *grpc_service.GetStorageItemsRequest) (*grpc_service.StorageItemsResponse, error) {
	var getDriveCommand commands.GetDriveCommand
	err := shared_kernel.MapToGrpc(&getDriveCommand, &req)

	dirs, err := f.getDriveHandler.Handle(getDriveCommand)
	if err != nil {
		f.logger.Error(err.Error())
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	response := grpc_service.StorageItemsResponse{}
	err = shared_kernel.MapToGrpc(&response.Directories, &dirs)
	if err != nil {
		f.logger.Error(err.Error())
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return &response, nil
}
