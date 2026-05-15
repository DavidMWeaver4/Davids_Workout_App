package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"

	"github.com/alexedwards/argon2id"
)

func GetBearerToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("Missing authorization token")
	}

	parts := strings.Fields(authHeader)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return "", fmt.Errorf("authorization header must be in format 'Bearer {token}'")
	}

	return strings.TrimSpace(parts[1]), nil
}

func MakeRefreshToken() (string, error) {
	token := make([]byte, 32)
	_, err := rand.Read(token)
	if err != nil {
		return "", err
	}
	string_token := hex.EncodeToString(token)

	return string_token, nil
}

func GetAPIKey(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("Missing authorization token")
	}

	parts := strings.Fields(authHeader)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "ApiKey") {
		return "", fmt.Errorf("authorization header must be in format 'ApiKey {token}'")
	}

	return strings.TrimSpace(parts[1]), nil
}

func HashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

func HashPassword(password string) (string, error) {
	hash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return "", err
	}
	return hash, nil
}

func CheckPasswordHash(password, hash string) (bool, error) {
	match, err := argon2id.ComparePasswordAndHash(password, hash)
	if err != nil {
		return false, err
	}
	return match, nil
}
