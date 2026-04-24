package domain

import "github.com/golang-jwt/jwt/v5"

type ClaimsDto struct {
	Username string `json:"username"`
	UserID   string `json:"user_id"`
	jwt.RegisteredClaims
}
