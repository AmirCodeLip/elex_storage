package pgx_repositories

import (
	"elex_storage/file_storage/internal/core_utils"
	"elex_storage/file_storage/internal/domain/entities"
	"elex_storage/file_storage/internal/test"
	pkg_entities "elex_storage/pkg/shared_kernel/entities"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFileMetadataRepository(t *testing.T) {
	appTest := test.InjectBase(t)
	fileRepository := CreateFileRepository(
		appTest.Logger,
		appTest.Db,
		appTest.Cfg,
	)

	// Prepare data for test
	txtFile := []byte("this is simple text file for test the project")
	checksum := core_utils.GetContentHash(txtFile)

	tests := []struct {
		name       string
		fileEntity entities.FileEntity
		wantErr    bool
	}{
		{
			name: "valid",
			fileEntity: entities.FileEntity{
				Id:          uuid.New(),
				CreatedAt:   time.Now(),
				ContentType: pkg_entities.Unown,
				Checksum:    checksum,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create new transaction
			tx, err := fileRepository.BeginTransaction()
			require.NoError(t, err)

			// Insert file into db
			insertErr := fileRepository.Insert(&tc.fileEntity, tx)
			assert.Error(t, insertErr)

			// Everything succeeded, commit the database transaction
			require.NoError(t, tx.Commit())
		})
	}

	// Check all files are exist
	files, err := fileRepository.GetAll()
	require.NoError(t, err)

	require.NotEmpty(t, files)

	appTest.App.Done()
}
