package grpc_server

import (
	"context"
	"elex_storage/file_metadata/internal/use_case/cqrs/command_handlers"
	"elex_storage/file_metadata/internal/use_case/cqrs/commands"
	"elex_storage/file_metadata/internal/use_case/dtos"
	"elex_storage/pkg/logger"
	"elex_storage/pkg/shared_kernel"
	"elex_storage/pkg/shared_kernel/grpc_service"
	"errors"
	"fmt"

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

func NewFileMetadataService(
	getDirectoriesHandler *command_handlers.GetDirectoriesHandler,
	createDirectoryHandler *command_handlers.CreateDirectoryHandler,
	getDriveHandler *command_handlers.GetDriveHandler,
	logger logger.Logger) grpc_service.FileMetadataServiceServer {
	return &FileMetadataService{getDirectoriesHandler: getDirectoriesHandler,
		createDirectoryHandler: createDirectoryHandler,
		getDriveHandler:        getDriveHandler,
		logger:                 logger}
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
	// Map request to command
	var getDriveCommand commands.GetDriveCommand
	if err := shared_kernel.MapToGrpc(&getDriveCommand, req); err != nil {
		f.logger.Error("failed to map request to command", "error", err)
		return nil, status.Error(codes.InvalidArgument, "invalid request format")
	}

	// Handle the command
	dirs, files, err := f.getDriveHandler.Handle(getDriveCommand)
	if err != nil {
		f.logger.Error("failed to get storage items", "error", err)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// Build response
	response := &grpc_service.StorageItemsResponse{}

	// Map directories
	if err := f.mapDirectories(&response.Directories, dirs); err != nil {
		return nil, err
	}

	// Map files
	if err := f.mapFiles(&response.Files, files); err != nil {
		return nil, err
	}

	return response, nil
}

func (f *FileMetadataService) mapDirectories(target interface{}, dirs interface{}) error {
	if err := shared_kernel.MapToGrpc(target, dirs); err != nil {
		f.logger.Error("failed to map directories", "error", err)
		return status.Error(codes.InvalidArgument, "failed to process directories")
	}
	return nil
}

func (f *FileMetadataService) mapFiles(target *[]*grpc_service.FileInfo, files []dtos.FileDto) error {
	mappedFiles := make([]*grpc_service.FileInfo, 0, len(files))

	for i, file := range files {
		dst := &grpc_service.FileInfo{}

		if err := shared_kernel.MapToGrpc(dst, &file); err != nil {
			f.logger.Error("failed to map file", "index", i, "error", err)
			return status.Error(codes.InvalidArgument, fmt.Sprintf("failed to process file at index %d", i))
		}

		dst.CreatedAt = shared_kernel.TimeToUnix(file.CreatedAt)
		dst.UpdateAt = shared_kernel.TimeToUnix(file.UpdateAt)
		mappedFiles = append(mappedFiles, dst)
	}

	*target = mappedFiles
	return nil
}
