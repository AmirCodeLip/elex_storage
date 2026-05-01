package repositories

import "elex_storage/file_metadata/internal/domain/entities"

type DirectoryMetadataRepository interface {
	Insert(directoryMetadataEntity *entities.DirectoryMetadataEntity) error
	GetRoot() (*entities.DirectoryMetadataEntity, error)
	GetDirectories() (*[]entities.DirectoryMetadataEntity, error)
	SetDirectoryId(fileMetadataEntity *entities.FileMetadataEntity) error
}
