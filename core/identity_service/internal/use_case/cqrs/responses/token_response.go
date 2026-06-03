package responses

import "github.com/google/uuid"

type TokenResponse struct {
	UserId       uuid.UUID
	AccessToken  string
	RefreshToken string
}
