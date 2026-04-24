package repositories

import "elex_storage/file_metadata/internal/domain/entities"

type FileMetadataRepository interface {
	Insert(fileMetadataEntity entities.FileMetadataEntity) error
	GetFiles() (*[]entities.FileMetadataEntity, error)
}
