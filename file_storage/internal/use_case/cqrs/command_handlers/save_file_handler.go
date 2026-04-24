package command_handlers

import (
	"elex_storage/file_storage/internal/core_utils"
	"elex_storage/file_storage/internal/domain"
	"elex_storage/file_storage/internal/domain/entities"
	"elex_storage/file_storage/internal/domain/messaging/publishers"
	"elex_storage/file_storage/internal/domain/repositories"
	"elex_storage/file_storage/internal/use_case/cqrs"
	"elex_storage/file_storage/internal/use_case/cqrs/commands"
	"elex_storage/pkg/logger"
	"elex_storage/pkg/shared_kernel/event_models"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
)

type SaveFileHandler struct {
	logger           logger.Logger
	config           *domain.ConfigEnv
	fileRepository   repositories.FileRepository
	pathUtil         *core_utils.PathUtil
	storagePublisher publishers.StoragePublisher
}

func NewSaveFileHandler(logger logger.Logger, config *domain.ConfigEnv,
	fileRepository repositories.FileRepository, pathUtil *core_utils.PathUtil,
	storagePublisher publishers.StoragePublisher) *SaveFileHandler {
	return &SaveFileHandler{logger, config, fileRepository, pathUtil, storagePublisher}
}

func (u *SaveFileHandler) Handle(cmd commands.SaveFileCommand) error {
	newFileId := uuid.New()

	// ToDo Amir codelip encrypt file and check it's content type
	contentType := u.pathUtil.GetContentType(cmd.Data)

	// First step: Insert the file into the database.
	fileEntity := entities.FileEntity{Id: newFileId, Name: cmd.Name, CreatedAt: time.Now(), ContentType: contentType}
	insertErr := u.fileRepository.Insert(fileEntity)
	if insertErr != nil {
		u.logger.Error(insertErr.Error())
		return insertErr
	}

	// Second step: Get the file path and check whether it exists or not.
	path, fullPath, _ := u.pathUtil.GetPath(fileEntity)
	_, dirErr := os.Stat(path)
	if dirErr != nil {
		// Create a directory if it does not exist
		if os.IsNotExist(dirErr) {
			err := os.MkdirAll(path, 0755)
			if err != nil {
				fmt.Println(err)
				return cqrs.InvalidFileErr
			}

		} else {
			u.logger.Error(dirErr.Error())
			return cqrs.InvalidFileErr
		}
	}

	// Third step: Save the file to the disk.
	errSaveFile := os.WriteFile(fullPath, *cmd.Data, 0644)
	if errSaveFile != nil {
		u.logger.Error(errSaveFile.Error())
		return cqrs.InvalidFileErr
	}

	u.storagePublisher.PublishFileCreated(event_models.FileCreated{
		Id:          newFileId,
		Name:        cmd.Name,
		ContentType: contentType,
		Size:        len(*cmd.Data),
		Drive:       u.config.DriveName,
	})

	return nil
}
