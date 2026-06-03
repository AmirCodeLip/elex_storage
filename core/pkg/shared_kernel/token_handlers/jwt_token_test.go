package token_handlers

import (
	"fmt"
	"testing"
	"time"
)

func TestGenerateKey(t *testing.T) {
	jwtToken := JwtToken{}
	privateKey, publickKey, err := jwtToken.GenerateKeyPair()
	if err != nil {
		t.Error(err.Error())
	} else {
		// Export keys to PEM format
		privateKeyPEM := jwtToken.ExportPrivateKey(privateKey)
		publicKeyPEM, err := jwtToken.ExportPublicKey(publickKey)
		if err != nil {
			t.Error(err.Error())
		}
		fmt.Printf("private key is %s", privateKeyPEM)
		fmt.Printf("public key is %s", publicKeyPEM)
	}
}

func createJwtToken() (*JwtToken, error) {
	jwtToken := JwtToken{accessTokenDuration: time.Hour * 3, refreshTokenDuration: time.Hour * 3}

	// Generate RSA key pair
	privateKey, publicKey, err := jwtToken.GenerateKeyPair()
	if err != nil {
		return nil, err
	}
	jwtToken.SetTokens(privateKey, publicKey)
	return &jwtToken, nil
}

func TestToken(t *testing.T) {
	jwtToken, err := createJwtToken()
	if err != nil {
		t.Error(err.Error())
		return
	}
	tokenPair, err := jwtToken.GenerateTokenPair("user123", "admin")
	if err != nil {
		t.Error(err.Error())
		return
	}
	fmt.Printf("Access Token: %s\n\n", tokenPair.AccessToken)
	fmt.Printf("Refresh Token: %s\n\n", tokenPair.RefreshToken)
	fmt.Printf("Expires At: %v\n", time.Unix(tokenPair.ExpiresAt, 0))
}

func TestValidateToken(t *testing.T) {
	jwtToken, err := createJwtToken()
	if err != nil {
		t.Error(err.Error())
		return
	}
	tokenPair, err := jwtToken.GenerateTokenPair("user123", "admin")
	if err != nil {
		t.Error(err.Error())
		return
	}
	claims, _ := jwtToken.ValidateAccessToken(tokenPair.AccessToken)
	fmt.Println(claims)
}
