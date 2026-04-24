package pgx_repositories

import (
	"elex_storage/file_storage/internal/domain"
	"elex_storage/file_storage/internal/domain/entities"
	"elex_storage/file_storage/internal/domain/repositories"
	"elex_storage/pkg/logger"

	"github.com/jmoiron/sqlx"
)

type FileRepository struct {
	logger    logger.Logger
	db        *sqlx.DB
	driveDisk string
}

func CreateFileRepository(logger logger.Logger, db *sqlx.DB, config *domain.ConfigEnv) repositories.FileRepository {
	return &FileRepository{logger: logger, db: db, driveDisk: config.DriveDisk}
}

func (fileRepository *FileRepository) Insert(fileEntity entities.FileEntity) error {
	_, insertErr := fileRepository.db.Exec(
		"INSERT INTO files (id, name, is_mounted, created_at, content_type) VALUES ($1, $2, $3, $4, $5)",
		fileEntity.Id,
		fileEntity.Name,
		false,
		fileEntity.CreatedAt,
		int(fileEntity.ContentType),
	)
	if insertErr != nil {
		fileRepository.logger.Error(insertErr.Error())
		return InsertFileErr
	}
	return nil
}
