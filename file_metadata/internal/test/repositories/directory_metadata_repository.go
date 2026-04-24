package repositories

import (
	"crypto/sha256"
	"elex_storage/pkg/logger"
	"encoding/hex"

	"github.com/google/uuid"

	"elex_storage/file_metadata/internal/domain/entities"
	"elex_storage/file_metadata/internal/domain/repositories"
)

type MockDirectoryMetadataRepository struct {
	logger      logger.Logger
	root        *entities.DirectoryMetadataEntity
	directories map[uuid.UUID]*entities.DirectoryMetadataEntity
}

func CreateMockDirectoryMetadataRepository(logger logger.Logger) repositories.DirectoryMetadataRepository {
	return &MockDirectoryMetadataRepository{logger: logger,
		directories: make(map[uuid.UUID]*entities.DirectoryMetadataEntity)}
}

func (repo *MockDirectoryMetadataRepository) Insert(directoryMetadataEntity *entities.DirectoryMetadataEntity) error {
	if directoryMetadataEntity.ParentId == nil {
		// Set root for the parent id
		root, _ := repo.GetRoot()
		directoryMetadataEntity.ParentId = &root.Id
	}
	directoryMetadataEntity.Id = uuid.New()
	directoryMetadataEntity.Hash = repo.Hash(directoryMetadataEntity)
	repo.directories[directoryMetadataEntity.Id] = directoryMetadataEntity
	return nil
}

func (repo *MockDirectoryMetadataRepository) GetDirectories() (*[]entities.DirectoryMetadataEntity, error) {
	dirs := make([]entities.DirectoryMetadataEntity, 0, len(repo.directories))
	for _, dir := range repo.directories {
		dirs = append(dirs, *dir)
	}
	return &dirs, nil
}

func (repo *MockDirectoryMetadataRepository) GetRoot() (*entities.DirectoryMetadataEntity, error) {
	if repo.root != nil {
		return repo.root, nil
	} else {
		root := entities.DirectoryMetadataEntity{Name: "~"}
		// Create new root.
		root.Id = uuid.New()
		root.Hash = repo.Hash(&root)
		root.ParentId = nil
		repo.root = &root
		repo.directories[root.Id] = &root
		return repo.root, nil
	}
}

func (repo *MockDirectoryMetadataRepository) Hash(directoryMetadataEntity *entities.DirectoryMetadataEntity) string {
	hasher := sha256.New()
	hasher.Write([]byte(directoryMetadataEntity.Name))
	hash := hasher.Sum(nil)
	if directoryMetadataEntity.ParentId != nil {
		hasher.Write([]byte(directoryMetadataEntity.ParentId.String()))
		hash = hasher.Sum(nil)
	}
	return hex.EncodeToString(hash)
}
