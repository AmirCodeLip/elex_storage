package pgx_repositories

import (
	"database/sql"
	"elex_storage/file_storage/internal/domain/entities"
	"elex_storage/file_storage/internal/domain/repositories"
	"elex_storage/pkg/logger"
	"elex_storage/pkg/shared_kernel/models"

	"github.com/jmoiron/sqlx"
)

type FileRepository struct {
	logger    logger.Logger
	db        *sqlx.DB
	driveDisk string
}

func CreateFileRepository(logger logger.Logger, db *sqlx.DB, config *models.ConfigEnv) repositories.FileRepository {
	return &FileRepository{logger: logger, db: db, driveDisk: config.DriveDisk}
}

func (fileRepository *FileRepository) BeginTransaction() (*sqlx.Tx, error) {
	tx, err := fileRepository.db.Beginx()
	return tx, err
}

func (fileRepository *FileRepository) Insert(fileEntity *entities.FileEntity, tx *sqlx.Tx) error {
	_, insertErr := tx.NamedExec(
		`INSERT INTO files (id ,name ,is_mounted ,created_at ,content_type ,checksum) VALUES 
			(:id ,:name ,false ,:created_at ,:content_type ,:checksum)`,
		fileEntity,
	)
	if insertErr != nil {
		fileRepository.logger.Error(insertErr.Error())
		return InsertFileErr
	}
	return nil
}

func (r *FileRepository) GetByHash(checksum string) (*entities.FileEntity, error) {
	var f entities.FileEntity

	err := r.db.Get(&f, `SELECT * FROM Files WHERE checksum = $1`, checksum)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &f, nil
}
