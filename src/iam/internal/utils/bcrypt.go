package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"

	"github.com/sujithsss01/email_server/iam/internal/apperrors"
	"golang.org/x/crypto/bcrypt"
)

const BcryptCost = 12

func HashPassword(password string) (string, error) {
	passwordBytes := []byte(password)
	var emtpyresult string
	if len(passwordBytes) > 72 {
		return emtpyresult, apperrors.ErrPasswordTooLong
	}
	hashedPassword, err := bcrypt.GenerateFromPassword(passwordBytes, BcryptCost)
	if err != nil {
		return emtpyresult, err
	}
	return string(hashedPassword), nil
}

func VerifyPasswordAndHash(password string, hash string) error {
	passBytes := []byte(password)
	if len(passBytes) > 72 {
		return apperrors.ErrInvalidCredentials
	}
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return apperrors.ErrInvalidCredentials
	}
	return nil
}

func GenerateRefreshToken() (string, error) {
	tokenBytes := make([]byte, 32)
	_, err := rand.Read(tokenBytes)
	if err != nil {
		return "", apperrors.ErrUnexpectedInternalFailure
	}
	return base64.RawURLEncoding.EncodeToString(tokenBytes), nil
}

func HashRefreshToken(refreshToken string) string {
	hash := sha256.Sum256([]byte(refreshToken))
	return hex.EncodeToString(hash[:])
}
