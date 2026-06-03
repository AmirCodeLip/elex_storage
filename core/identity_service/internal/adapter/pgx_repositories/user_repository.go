package pgx_repositories

import (
	"database/sql"
	"elex_storage/identity_service/internal/domain/entities"
	domain_errors "elex_storage/identity_service/internal/domain/errors"
	"elex_storage/identity_service/internal/domain/repositories"
	"elex_storage/pkg/logger"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	logger    logger.Logger
	db        *sqlx.DB
	driveDisk string
}

func CreateFileRepository(logger logger.Logger, db *sqlx.DB) repositories.UserRepository {
	return &UserRepository{logger: logger, db: db}
}

func (repo *UserRepository) Insert(userEntity entities.UserEntity) error {
	if *userEntity.Phone == "" {
		userEntity.Phone = nil
	}
	if *userEntity.Email == "" {
		userEntity.Email = nil
	}
	_, insertErr := repo.db.NamedExec(`INSERT INTO users(id, email, phone, password_hash) 
		VALUES (:id, :email, :phone, :password_hash)`, userEntity)
	if insertErr != nil {
		var pgErr *pgconn.PgError
		if errors.As(insertErr, &pgErr) {
			if pgErr.Code == "23505" {
				return domain_errors.ErrUserExists
			}
		}
		repo.logger.Error(insertErr.Error())
		return domain_errors.ErrInsertUser
	}
	return nil
}

func (repo *UserRepository) GetUserByPhone(phone string) (*entities.UserEntity, error) {
	var userEntity entities.UserEntity
	getErr := repo.db.Get(&userEntity, `SELECT * FROM public.users WHERE phone = $1`, phone)
	if getErr != nil {
		if errors.Is(getErr, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, getErr
	}
	return &userEntity, nil
}

func (repo *UserRepository) GetUserByEmail(email string) (*entities.UserEntity, error) {
	var userEntity entities.UserEntity
	getErr := repo.db.Get(&userEntity, `SELECT * FROM public.users WHERE email = $1`, email)
	if getErr != nil {
		if errors.Is(getErr, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, getErr
	}
	return &userEntity, nil
}
