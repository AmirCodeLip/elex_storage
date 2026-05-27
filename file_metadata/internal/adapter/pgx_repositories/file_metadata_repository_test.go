package pgx_repositories

import (
	"elex_storage/file_metadata/internal/domain/entities"
	"elex_storage/file_metadata/internal/test"
	"elex_storage/pkg/shared_kernel/models"
	"errors"
	"strconv"
	"testing"

	"math/rand"

	"github.com/google/uuid"
)

func TestFileMetadataRepository(t *testing.T) {
	_, db, logger := test.InjectBase(t)
	directoryMetadataRepository := CreateDirectoryMetadataRepository(logger, db)
	fileMetadataRepository := CreateFileMetadataRepository(logger, db)

	root, err := directoryMetadataRepository.GetRoot()
	if err != nil {
		t.Fatalf("failed to get root directory: %v", err)
	}

	tests := []struct {
		name    string
		file    entities.FileMetadataEntity
		wantErr bool
	}{
		{
			name: "insert valid file",
			file: entities.FileMetadataEntity{
				Id:            uuid.New(),
				Name:          "test_file_" + strconv.Itoa(rand.Intn(10000)),
				Size:          10,
				DirectoryId:   &root.Id,
				StorageId:     uuid.New(),
				Checksum:      "VALID CHECKSUM",
				FileExtension: ".txt",
				Drive:         "test_drive",
				EncryptionKey: []byte{92, 100, 54, 60},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			tx, err := fileMetadataRepository.BeginTransaction()
			defer tx.Rollback()
			if err != nil {
				t.Errorf("Insert() error = %v, wantErr %v", err, tt.wantErr)
			}

			err = fileMetadataRepository.Insert(tt.file, tx)
			if (err != nil) != tt.wantErr {
				var commonErr *models.CommonError
				if errors.As(err, &commonErr) {
					/// File is uploaded before just check exist
					logger.Error("File is uploaded befor")
					return
				} else {
					logger.Error(err.Error())
					t.Errorf("Insert() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
			}
			tx.Commit()

			if !tt.wantErr {
				// Verify the file was actually inserted
				files, err := fileMetadataRepository.GetFiles(root.Id)
				if err != nil {
					t.Fatalf("failed to fetch files: %v", err)
				}

				found := false
				for _, f := range *files {
					if f.Id == tt.file.Id {
						found = true
						if f.Name != tt.file.Name {
							t.Errorf("file name = %v, want %v", f.Name, tt.file.Name)
						}
						// Add more field validations...
						break
					}
				}

				if !found {
					t.Error("inserted file not found in repository")
				}
			}
		})
	}
}
