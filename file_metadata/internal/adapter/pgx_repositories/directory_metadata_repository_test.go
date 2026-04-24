package pgx_repositories

import (
	"elex_storage/file_metadata/internal/domain/entities"
	"elex_storage/file_metadata/internal/test"
	"elex_storage/pkg/shared_kernel/models"
	"errors"
	"fmt"
	"testing"

	"github.com/google/uuid"
)

func TestDirectoryMetadataRepository(t *testing.T) {
	_, db, logger := test.InjectBase(t)
	directoryMetadataRepository := CreateDirectoryMetadataRepository(logger, db)
	// Insert directories
	insertErr1 := directoryMetadataRepository.Insert(&entities.DirectoryMetadataEntity{
		Id: uuid.New(), Name: "test_dir1", Size: 10,
	})
	if insertErr1 != nil {
		var commonErr *models.CommonError
		if errors.As(insertErr1, &commonErr) {
			logger.Error("Can't create directory", commonErr.Error())
		} else {
			logger.Error(insertErr1.Error())
			return
		}
	}
	insertErr2 := directoryMetadataRepository.Insert(&entities.DirectoryMetadataEntity{
		Id: uuid.New(), Name: "test_dir2", Size: 10,
	})
	if insertErr2 != nil {
		var commonErr *models.CommonError
		if errors.As(insertErr2, &commonErr) {
			logger.Error("Can't create directory", commonErr.Error())
		} else {
			logger.Error(insertErr2.Error())
		}
	}
	dirs, err := directoryMetadataRepository.GetDirectories()
	if err != nil {
		logger.Error("failed to fetch directories: %v", err)
	}
	if dirs == nil || len(*dirs) < 1 {
		logger.Error("expected one or more directories, got %d", len(*dirs))
	}
}

func TestGetDirectoryMetadataRepository(t *testing.T) {
	_, db, logger := test.InjectBase(t)
	directoryMetadataRepository := CreateDirectoryMetadataRepository(logger, db)
	dirs, getErr1 := directoryMetadataRepository.GetDirectories()
	if getErr1 != nil {
		logger.Error(getErr1.Error())
		return
	}
	for _, dir := range *dirs {
		fmt.Println(dir.Name)
	}
}
