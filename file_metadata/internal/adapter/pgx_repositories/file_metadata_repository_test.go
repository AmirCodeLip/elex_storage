package pgx_repositories

import (
	"elex_storage/file_metadata/internal/domain/entities"
	"elex_storage/file_metadata/internal/test"
	"testing"

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
				Name:          "test_file_1",
				Size:          10,
				DirectoryId:   &root.Id,
				StorageId:     uuid.New(),
				Checksum:      "VALID CHECKSUM",
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
				t.Errorf("Insert() error = %v, wantErr %v", err, tt.wantErr)
			}
			tx.Commit()

			if !tt.wantErr {
				// Verify the file was actually inserted
				files, err := fileMetadataRepository.GetFiles()
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
