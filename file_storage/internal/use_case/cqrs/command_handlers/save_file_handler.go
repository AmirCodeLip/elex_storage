package command_handlers

import (
	"elex_storage/file_storage/internal/core_utils"
	"elex_storage/file_storage/internal/domain/entities"
	"elex_storage/file_storage/internal/domain/messaging/publishers"
	"elex_storage/file_storage/internal/domain/repositories"
	"elex_storage/file_storage/internal/use_case/cqrs"
	"elex_storage/file_storage/internal/use_case/cqrs/commands"
	"elex_storage/pkg/logger"
	"elex_storage/pkg/shared_kernel/event_models"
	"elex_storage/pkg/shared_kernel/models"
	"os"
	"time"

	"github.com/google/uuid"
)

type SaveFileHandler struct {
	logger           logger.Logger
	config           *models.ConfigEnv2
	fileRepository   repositories.FileRepository
	pathUtil         *core_utils.PathUtil
	storagePublisher publishers.StoragePublisher
}

func NewSaveFileHandler(logger logger.Logger, config *models.ConfigEnv2,
	fileRepository repositories.FileRepository, pathUtil *core_utils.PathUtil,
	storagePublisher publishers.StoragePublisher) *SaveFileHandler {
	return &SaveFileHandler{logger, config, fileRepository, pathUtil, storagePublisher}
}

func (u *SaveFileHandler) Handle(cmd commands.SaveFileCommand) error {
	// Step 0: Guard clause for nil pointer
	if cmd.Data == nil {
		u.logger.Error("cmd.Data is nil")
		return cqrs.SaveFileErr(cmd.Name)
	}

	checksum := core_utils.GetContentHash(*cmd.Data)
	var key []byte = nil

	// Step 1: Check if file with same hash exists
	fileEntity, err := u.fileRepository.GetByHash(checksum)
	if err != nil {
		u.logger.Error(err.Error())
		return cqrs.SaveFileErr(cmd.Name)
	}

	if fileEntity != nil {
		// File exists in storage.
	} else {
		newFileId := uuid.New()
		contentType := u.pathUtil.GetContentType(cmd.Data)

		key, err = core_utils.GenerateNewKey()
		if err != nil {
			u.logger.Error(err.Error())
			return cqrs.SaveFileErr(cmd.Name)
		}

		fileEntity = &entities.FileEntity{
			Id:          newFileId,
			Name:        cmd.Name,
			CreatedAt:   time.Now(),
			ContentType: contentType,
			Checksum:    checksum,
		}

		// Step 2: Get the file path and create directory if needed (Moved BEFORE DB insert to prevent dirty DB state)
		path, fullPath := u.pathUtil.GetPath(*fileEntity)

		if _, dirErr := os.Stat(path); os.IsNotExist(dirErr) {
			if mkdirErr := os.MkdirAll(path, 0755); mkdirErr != nil {
				u.logger.Error(mkdirErr.Error())
				return cqrs.SavePathErr(fullPath)
			}
		} else if dirErr != nil {
			u.logger.Error(dirErr.Error())
			return cqrs.SavePathErr(fullPath)
		}

		// Step 3: Encrypt content of the file.
		content, encErr := core_utils.EncryptContent(*cmd.Data, key)
		if encErr != nil {
			u.logger.Error(encErr.Error())
			return cqrs.SaveFileErr(cmd.Name)
		}

		// Step 4: Start the database transaction
		tx, err := u.fileRepository.BeginTransaction()
		if err != nil {
			u.logger.Error(err.Error())
			return cqrs.SaveFileErr(fullPath)
		}

		// It guarantees the DB is rolled back if a panic occurs or we return early.
		defer tx.Rollback()

		// Step 5: Insert the file into the database.
		err = u.fileRepository.Insert(fileEntity, tx)
		if err != nil {
			u.logger.Error(err.Error())
			return err
		}

		// Step 6: Save the file to the disk.
		errSaveFile := os.WriteFile(fullPath, content, 0644)
		if errSaveFile != nil {
			u.logger.Error(errSaveFile.Error())
			return cqrs.SavePathErr(fullPath)
		}

		// 7. Everything succeeded, commit the database transaction
		if commitErr := tx.Commit(); commitErr != nil {
			u.logger.Error(commitErr.Error())
			// Edge case: If DB commit fails AFTER file is written, delete the orphaned file
			os.Remove(fullPath)
			return cqrs.SaveFileErr(fullPath)
		}

	}

	// Step 6: Publish Event
	u.storagePublisher.PublishFileCreated(event_models.FileCreated{
		Id:          fileEntity.Id,
		Name:        cmd.Name,
		ContentType: fileEntity.ContentType,
		Size:        len(*cmd.Data),
		Drive:       u.config.ServiceName,
		Checksum:    checksum,
		Key:         key, // Make sure 'key' isn't expected to be populated for existing files if left nil
	})

	return nil
}
