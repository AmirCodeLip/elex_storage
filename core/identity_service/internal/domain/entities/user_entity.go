package entities

import (
	"time"

	"github.com/google/uuid"
)

type UserEntity struct {
	Id           uuid.UUID `db:"id"`
	Email        *string   `db:"email"`
	PasswordHash string    `db:"password_hash"`
	Phone        *string   `db:"phone"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}
