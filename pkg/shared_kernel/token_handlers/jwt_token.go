package token_handlers

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"elex_storage/pkg/shared_kernel/token_handlers/models"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JwtToken struct {
	privateKey           *rsa.PrivateKey
	publicKey            *rsa.PublicKey
	accessTokenDuration  time.Duration
	refreshTokenDuration time.Duration
}

// func NewJwtToken(config *domain.ConfigEnv) *JwtToken {
// 	return &JwtToken{accessTokenDuration: config.AccessTokenDuration, refreshTokenDuration: config.RefreshTokenDuration}
// }

func NewJwtToken(accessTokenDuration time.Duration, refreshTokenDuration time.Duration) (*JwtToken, error) {
	jwtToken := JwtToken{accessTokenDuration: accessTokenDuration, refreshTokenDuration: refreshTokenDuration}
	return &jwtToken, nil
}

func (jwtToken *JwtToken) SetPublicKey() error {
	pwd, _ := os.Getwd()
	var exeDir = filepath.Dir(pwd)
	// 1. Read public key
	pathPublicKey := filepath.Join(exeDir, "jwt_identity_public.key")
	file2, err := os.Open(pathPublicKey)
	if err != nil {
		return errors.New(fmt.Sprintf("Failed to read key file: %s", pathPublicKey))
	}
	defer file2.Close()
	publicKeyData, err := io.ReadAll(file2)
	// 4. Parse public key
	publicKey, err := jwtToken.ParsePublicKey(publicKeyData)
	if err != nil {
		return errors.Join(errors.New("Invalid public key is provided"), err)
	}
	jwtToken.publicKey = publicKey
	return nil
}

func (jwtToken *JwtToken) SetPrivateKey() error {
	pwd, _ := os.Getwd()
	var exeDir = filepath.Dir(pwd)
	// 1. Read private key
	pathPrivateKey := filepath.Join(exeDir, "jwt_identity_private.key")
	file1, err := os.Open(pathPrivateKey)
	if err != nil {
		return errors.New(fmt.Sprintf("Failed to read key file: %s", pathPrivateKey))
	}
	defer file1.Close()
	privateKeyData, err := io.ReadAll(file1)
	if err != nil {
		return err
	}
	// 2. Parse private key
	privateKey, err := jwtToken.ParsePrivateKey(privateKeyData)
	if err != nil {
		return errors.Join(errors.New("Invalid private key is provided"), err)
	}
	jwtToken.privateKey = privateKey
	return nil
}

func (jwtToken *JwtToken) SetTokens(privateKey *rsa.PrivateKey, publicKey *rsa.PublicKey) {
	jwtToken.privateKey = privateKey
	jwtToken.publicKey = publicKey
}

func (jwtToken *JwtToken) GenerateKeyPair() (*rsa.PrivateKey, *rsa.PublicKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}
	return privateKey, &privateKey.PublicKey, nil
}

// GenerateTokenPair generates both access and refresh tokens
func (jwtToken *JwtToken) GenerateTokenPair(userID, role string) (*models.TokenPair, error) {
	// Generate access token
	accessToken, accessExp, err := jwtToken.generateAccessToken(userID, role)
	if err != nil {
		return nil, err
	}

	// Generate refresh token
	refreshToken, err := jwtToken.generateRefreshToken(userID)
	if err != nil {
		return nil, err
	}

	return &models.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    accessExp.Unix(),
	}, nil
}

// generateAccessToken generates an access token
func (jwtToken *JwtToken) generateAccessToken(userID, role string) (string, time.Time, error) {
	expirationTime := time.Now().Add(jwtToken.accessTokenDuration)

	claims := &models.Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   userID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tokenString, err := token.SignedString(jwtToken.privateKey)
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expirationTime, nil
}

// ExportPrivateKey exports private key to PEM format
func (jwtToken *JwtToken) ExportPrivateKey(privateKey *rsa.PrivateKey) []byte {
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	})
	return privateKeyPEM
}

// ExportPublicKey exports public key to PEM format
func (jwtToken *JwtToken) ExportPublicKey(publicKey *rsa.PublicKey) ([]byte, error) {
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return nil, err
	}
	publicKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: publicKeyBytes,
	})
	return publicKeyPEM, nil
}

func (jwtToken *JwtToken) ParsePrivateKey(data []byte) (*rsa.PrivateKey, error) {
	// Try to decode as PEM
	block, _ := pem.Decode(data)
	if block == nil {
		// Not PEM encoded, try parsing as raw DER
		return parsePrivateKey(data)
	}

	// Handle different PEM types
	switch block.Type {
	case "RSA PRIVATE KEY":
		return x509.ParsePKCS1PrivateKey(block.Bytes)
	case "PRIVATE KEY":
		key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		if rsaKey, ok := key.(*rsa.PrivateKey); ok {
			return rsaKey, nil
		}
		return nil, fmt.Errorf("not an RSA private key")
	case "ENCRYPTED PRIVATE KEY":
		return nil, fmt.Errorf("encrypted private keys require password")
	default:
		return nil, fmt.Errorf("unsupported key type: %s", block.Type)
	}
}

func (jwtToken *JwtToken) ParsePublicKey(data []byte) (*rsa.PublicKey, error) {
	// Decode PEM block
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, fmt.Errorf("failed to parse PEM block containing the key")
	}

	// Parse the public key
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		// Try parsing as PKCS1 public key
		return x509.ParsePKCS1PublicKey(block.Bytes)
	}

	// Type assert to RSA public key
	rsaPub, ok := pub.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("key is not an RSA public key")
	}

	return rsaPub, nil
}

func parsePrivateKey(der []byte) (*rsa.PrivateKey, error) {
	// Try PKCS1
	if key, err := x509.ParsePKCS1PrivateKey(der); err == nil {
		return key, nil
	}

	// Try PKCS8
	if key, err := x509.ParsePKCS8PrivateKey(der); err == nil {
		if rsaKey, ok := key.(*rsa.PrivateKey); ok {
			return rsaKey, nil
		}
	}

	return nil, fmt.Errorf("failed to parse private key")
}

// generateRefreshToken generates a refresh token
func (jwtToken *JwtToken) generateRefreshToken(userID string) (string, error) {
	expirationTime := time.Now().Add(jwtToken.refreshTokenDuration)

	claims := &models.Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   userID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tokenString, err := token.SignedString(jwtToken.privateKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateAccessToken validates an access token
func (jwtToken *JwtToken) ValidateAccessToken(tokenString string) (*models.Claims, error) {
	claims := &models.Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtToken.publicKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}

// ValidateRefreshToken validates a refresh token
func (jwtToken *JwtToken) ValidateRefreshToken(tokenString string) (*models.Claims, error) {
	claims := &models.Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtToken.publicKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid refresh token")
	}

	return claims, nil
}

// RefreshTokenPair refreshes both tokens using a valid refresh token
func (jwtToken *JwtToken) RefreshTokenPair(refreshToken string) (*models.TokenPair, error) {
	// Validate refresh token
	claims, err := jwtToken.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %v", err)
	}

	// Generate new token pair
	return jwtToken.GenerateTokenPair(claims.UserID, claims.Role)
}
