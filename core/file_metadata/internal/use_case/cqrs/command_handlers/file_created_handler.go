package command_handlers

import (
	"elex_storage/file_metadata/internal/domain/entities"
	"elex_storage/file_metadata/internal/domain/repositories"
	"elex_storage/file_metadata/internal/use_case/cqrs"
	"elex_storage/file_metadata/internal/use_case/cqrs/commands"
	"elex_storage/pkg/logger"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/jinzhu/copier"
)

type FileCreatedHandler struct {
	logger                      logger.Logger
	fileMetadataRepository      repositories.FileMetadataRepository
	directoryMetadataRepository repositories.DirectoryMetadataRepository
}

func NewFileCreatedHandler(logger logger.Logger,
	fileMetadataRepository repositories.FileMetadataRepository, directoryMetadataRepository repositories.DirectoryMetadataRepository) *FileCreatedHandler {
	return &FileCreatedHandler{logger, fileMetadataRepository, directoryMetadataRepository}
}

func (u *FileCreatedHandler) Handle(cmd commands.FileCreatedCommand) error {
	// Step 1: Map event into entity file
	var fileMetadataEntity entities.FileMetadataEntity
	if err := copier.Copy(&fileMetadataEntity, &cmd.FileEntity); err != nil {
		u.logger.Error("Failed to copy file entity: " + err.Error())
		return err
	}

	// Step 2: Assign the data
	fileMetadataEntity.Id = uuid.New()
	fileMetadataEntity.StorageId = cmd.FileEntity.Id

	// Extract the file extension and clean the file name
	ext := filepath.Ext(fileMetadataEntity.Name)
	fileMetadataEntity.Name = strings.TrimSuffix(fileMetadataEntity.Name, ext)
	fileMetadataEntity.FileExtension = ext

	// Step 3: Set directory (default to root if not provided)
	if err := u.directoryMetadataRepository.SetDirectoryId(&fileMetadataEntity); err != nil {
		u.logger.Error("Failed to set directory ID: " + err.Error())
		return err
	}

	// Step 4: Execute insert operation with transaction
	tx, err := u.fileMetadataRepository.BeginTransaction()
	if err != nil {
		u.logger.Error("Failed to begin transaction: " + err.Error())
		return cqrs.SaveFileMetadataErr()
	}

	// Ensure transaction is rolled back on error
	defer func() {
		if p := recover(); p != nil {
			_ = u.fileMetadataRepository.RollbackTransaction(tx)
			panic(p) // Re-panic after rollback
		}
	}()

	// Insert the file
	if err = u.fileMetadataRepository.Insert(fileMetadataEntity, tx); err != nil {
		// Rollback transaction on insert error
		if rbErr := u.fileMetadataRepository.RollbackTransaction(tx); rbErr != nil {
			u.logger.Error("Failed to rollback transaction: " + rbErr.Error())
		}

		if isCommonErr, err := logger.HandleCommonErr(err, u.logger); err != nil {
			if isCommonErr {
				return err
			}
			return cqrs.SaveFileMetadataErr()
		}
	}

	// Commit the transaction
	if err = u.fileMetadataRepository.CommitTransaction(tx); err != nil {
		u.logger.Error("Failed to commit transaction: " + err.Error())
		return cqrs.SaveFileMetadataErr()
	}
	return nil
}
