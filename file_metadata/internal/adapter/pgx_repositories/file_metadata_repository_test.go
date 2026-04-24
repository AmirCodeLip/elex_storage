package pgx_repositories

import (
	"elex_storage/file_metadata/internal/domain/entities"
	"elex_storage/file_metadata/internal/test"
	"elex_storage/pkg/shared_kernel/models"
	"errors"
	"testing"

	"github.com/google/uuid"
)

func TestFileMetadataRepository(t *testing.T) {
	_, db, logger := test.InjectBase(t)
	directoryMetadataRepository := CreateDirectoryMetadataRepository(logger, db)
	fileMetadataRepository := CreateFileMetadataRepository(logger, db)
	root, _ := directoryMetadataRepository.GetRoot()
	// Insert directories
	insertErr1 := fileMetadataRepository.Insert(entities.FileMetadataEntity{
		Id: uuid.New(), Name: "test_file_1", Size: 10,
		DirectoryId: &root.Id, StorageId: uuid.New(),
	})
	if insertErr1 != nil {
		var commonErr *models.CommonError
		if errors.As(insertErr1, &commonErr) {
			t.Logf("Can't create file %s \n", commonErr.Error())
		} else {
			t.Error(insertErr1.Error())
		}
	}
	files, err := fileMetadataRepository.GetFiles()
	if err != nil {
		t.Fatalf("failed to fetch files: %v", err)
	}
	if files == nil || len(*files) < 1 {
		t.Errorf("expected one or more files, got %d", len(*files))
	}
}
