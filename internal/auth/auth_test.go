package auth

import "testing"

func TestCheckPasswordHash(t *testing.T) {
	password := "testingpass123"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword returned error: %v", err)
	}

	if hash == password {
		t.Fatal("password was not hashed")
	}

	match, err := CheckPasswordHash(password, hash)
	if err != nil {
		t.Fatalf("CheckPasswordHash returned error: %v", err)
	}

	if !match {
		t.Fatal("expected password to match hash")
	}
}
func TestCheckPasswordHash_WrongPassword(t *testing.T) {
	hash, err := HashPassword("testingpass123")
	if err != nil {
		t.Fatalf("HashPassword returned error: %v", err)
	}

	match, err := CheckPasswordHash("wrong-password", hash)
	if err != nil {
		t.Fatalf("CheckPasswordHash returned error: %v", err)
	}

	if match {
		t.Fatal("expected password validation to fail")
	}
}
func TestCheckPasswordHash_EmptyPassword(t *testing.T) {
	hash, err := HashPassword("testingpass123")
	if err != nil {
		t.Fatalf("HashPassword returned error: %v", err)
	}

	match, err := CheckPasswordHash("", hash)
	if err != nil {
		t.Fatalf("CheckPasswordHash returned error: %v", err)
	}

	if match {
		t.Fatal("expected empty password to not match")
	}
}
func TestHashPasswordMakesUniqueHash(t *testing.T) {
	hash1, err := HashPassword("password123")
	if err != nil {
		t.Fatalf("error hasing password: %v", err)
	}

	hash2, err := HashPassword("password123")
	if err != nil {
		t.Fatalf("error hasing password: %v", err)
	}

	if hash1 == hash2 {
		t.Fatal("expected unique hashes")
	}
}
