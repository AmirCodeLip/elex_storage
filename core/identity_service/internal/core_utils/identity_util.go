package core_utils

import (
	"crypto/sha256"
	"elex_storage/pkg/shared_kernel/models"
	"encoding/hex"
)

type IdentityUtil struct {
	config *models.ConfigEnv
}

func NewIdentityUtil() *IdentityUtil {
	return &IdentityUtil{}
}

func (identityUtil *IdentityUtil) HashPassword(password string) (string, error) {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:]), nil
}

// CheckPasswordHash compares a plain text password with a hashed password
func (identityUtil *IdentityUtil) CheckPasswordHash(password, hash string) bool {
	hashPass, _ := identityUtil.HashPassword(password)
	return hashPass == hash
}
