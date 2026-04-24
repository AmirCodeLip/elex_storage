package repositories

import "elex_storage/file_storage/internal/domain/entities"

type FileRepository interface {
	Insert(fileEntity entities.FileEntity) error
}
