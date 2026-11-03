package utils

import (
	"crypto/rsa"
	"errors"
	"os"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AccessClaims struct {
	SessionId string `json:"sid"`
	jwt.RegisteredClaims
}

var privateKey *rsa.PrivateKey
var once sync.Once

var publicKey *rsa.PublicKey

func loadPrivateAndPublicKey() {
	privateKeyPem, err := os.ReadFile("./private.pem")
	if err != nil {
		panic("failed to read JWT private key")
	}
	privateKey, err = jwt.ParseRSAPrivateKeyFromPEM(privateKeyPem)
	if err != nil {
		panic("failed to create private key")
	}

	publicKeyPem, err := os.ReadFile("./public.pem")
	if err != nil {
		panic("failed to read JWT public key")
	}
	publicKey, err = jwt.ParseRSAPublicKeyFromPEM(publicKeyPem)
	if err != nil {
		panic("failed to create public key")
	}
}

func GenerateJWT(userId string, sessionId string) (string, error) {
	once.Do(loadPrivateAndPublicKey)
	now := time.Now()
	claims := AccessClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userId,
			Issuer:    "iam-service",
			Audience:  jwt.ClaimStrings{"email-api"},
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(30 * time.Minute)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["kid"] = "iam-key-1"
	return token.SignedString(privateKey)
}

func ValidateJwt(tokenString string) (*AccessClaims, error) {
	once.Do(loadPrivateAndPublicKey)
	claim := &AccessClaims{}
	token, err := jwt.ParseWithClaims(
		tokenString,
		claim,
		func(token *jwt.Token) (interface{}, error) {
			if token.Method != jwt.SigningMethodRS256 {
				return nil, errors.New("unexpected signing algorithm")
			}
			kid, ok := token.Header["kid"].(string)
			if !ok || kid != "iam-key-1" {
				return nil, errors.New("unknown signing key")
			}

			return publicKey, nil
		},
		jwt.WithIssuer("iam-service"),
		jwt.WithAudience("email-api"),
		jwt.WithExpirationRequired(),
		jwt.WithValidMethods([]string{"RS256"}),
	)
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claim, nil
}
