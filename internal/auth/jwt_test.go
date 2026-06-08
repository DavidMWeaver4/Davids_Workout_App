package auth

import (
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestJWTCreation(t *testing.T) {
	userID := uuid.New()
	secret := "testSecret"

	token, err := MakeJWT(userID, secret, time.Hour)
	if err != nil {
		t.Fatalf("MakeJWT returned error: %v", err)
	}

	gotID, err := ValidateJWT(token, secret)
	if err != nil {
		t.Fatalf("ValidateJWT returned error: %v", err)
	}

	if gotID != userID {
		t.Fatalf("expected userID %v, got %v", userID, gotID)
	}
}
func TestJWTWrongSecret(t *testing.T) {
	userID := uuid.New()

	token, err := MakeJWT(userID, "correctSecret", time.Hour)
	if err != nil {
		t.Fatalf("MakeJWT returned error: %v", err)
	}

	_, err = ValidateJWT(token, "wrong-secret")

	if err == nil {
		t.Fatal("expected validation to fail")
	}
}

func TestGetBearerToken(t *testing.T) {
	headers := http.Header{}
	headers.Set("Authorization", "Bearer abc123")

	token, err := GetBearerToken(headers)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if token != "abc123" {
		t.Fatalf("expected token abc123, got %s", token)
	}
}

func TestGetBearerToken_InvalidFormat(t *testing.T) {
	headers := http.Header{}
	headers.Set("Authorization", "abc123")

	_, err := GetBearerToken(headers)

	if err == nil {
		t.Fatal("expected error")
	}
}

func TestJWTExpired(t *testing.T) {
	userID := uuid.New()

	token, err := MakeJWT(userID, "testSecret", -time.Hour)
	if err != nil {
		t.Fatalf("MakeJWT returned error: %v", err)
	}

	_, err = ValidateJWT(token, "testSecret")

	if err == nil {
		t.Fatal("expected expired token error")
	}
}
