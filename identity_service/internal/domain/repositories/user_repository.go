package repositories

import "elex_storage/identity_service/internal/domain/entities"

type UserRepository interface {
	Insert(userEntity entities.UserEntity) error
	GetUserByPhone(phone string) (*entities.UserEntity, error)
	GetUserByEmail(email string) (*entities.UserEntity, error)
}
