package repositories

import (
	"elex_storage/file_metadata/internal/domain/entities"

	"github.com/google/uuid"
)

type DirectoryMetadataRepository interface {
	Insert(directoryMetadataEntity *entities.DirectoryMetadataEntity) error
	GetRoot() (*entities.DirectoryMetadataEntity, error)
	GetDirectories(parentId uuid.UUID) (*[]entities.DirectoryMetadataEntity, error)
	SetDirectoryId(fileMetadataEntity *entities.FileMetadataEntity) error
}
