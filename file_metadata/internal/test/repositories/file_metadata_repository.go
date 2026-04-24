package repositories

import (
	"crypto/sha256"
	"elex_storage/pkg/logger"
	"encoding/hex"
	"errors"
	"time"

	"elex_storage/file_metadata/internal/domain/entities"
	domain_errors "elex_storage/file_metadata/internal/domain/errors"
	"elex_storage/file_metadata/internal/domain/repositories"
)

type MockFileMetadataRepository struct {
	logger logger.Logger
	root   *entities.DirectoryMetadataEntity
	files  map[string]*entities.FileMetadataEntity
}

func CreateMockFileMetadataRepository(logger logger.Logger) repositories.FileMetadataRepository {
	return &MockFileMetadataRepository{
		logger: logger,
		files:  make(map[string]*entities.FileMetadataEntity)}
}

func (repo *MockFileMetadataRepository) Insert(fileMetadataEntity entities.FileMetadataEntity) error {
	if fileMetadataEntity.DirectoryId == nil {
		return errors.New("DirectoryId can't be null")
	}
	fileMetadataEntity.Hash = repo.Hash(fileMetadataEntity)
	if repo.files[fileMetadataEntity.Hash] != nil {
		return domain_errors.ErrInsertFileMetadata
	}
	fileMetadataEntity.CreatedAt = time.Now()
	fileMetadataEntity.UpdateAt = time.Now()
	repo.files[fileMetadataEntity.Hash] = &fileMetadataEntity
	return nil
}

func (repo *MockFileMetadataRepository) GetFiles() (*[]entities.FileMetadataEntity, error) {
	files := make([]entities.FileMetadataEntity, 0, len(repo.files))
	for _, dir := range repo.files {
		files = append(files, *dir)
	}
	return &files, nil
}

func (repo *MockFileMetadataRepository) Hash(fileMetadataEntity entities.FileMetadataEntity) string {
	hasher := sha256.New()
	hasher.Write([]byte(fileMetadataEntity.DirectoryId.String()))
	hasher.Write([]byte(fileMetadataEntity.Name))
	hasher.Write([]byte(fileMetadataEntity.FileExtension))
	hash := hasher.Sum(nil)
	return hex.EncodeToString(hash)
}
