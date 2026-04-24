package event_models

import (
	"github.com/google/uuid"
)

type UserRegisterd struct {
	Id    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Email *string   `json:"email"`
	Phone *string   `json:"phone"`
}
