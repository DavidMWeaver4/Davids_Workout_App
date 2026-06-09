package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/DavidMWeaver4/Davids_Workout_App/internal/auth"
	"github.com/DavidMWeaver4/Davids_Workout_App/internal/database"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

func TestAuthorizeWorkoutSession_Success(t *testing.T) {
	t.Skip("not yet implemented")
	ctx := context.Background()

	cfg := newTestAPIConfig(t)

	user := createTestUser(t, cfg)
	session := createTestWorkoutSession(t, cfg, user.ID)

	gotSession, err := cfg.authorizeWorkoutSession(ctx, session.ID, user.ID)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if gotSession.ID != session.ID {
		t.Fatalf("expected session ID %v, got %v", session.ID, gotSession.ID)
	}

	if gotSession.UserID != user.ID {
		t.Fatalf("expected user ID %v, got %v", user.ID, gotSession.UserID)
	}
}

func newTestAPIConfig(t *testing.T) *apiConfig {
	t.Helper()
	godotenv.Load(".env")
	db, err := sql.Open("postgres", os.Getenv("TEST_DB_URL"))
	if err != nil {
		t.Fatalf("failed to connect to test db: %v", err)
	}
	if !strings.Contains(os.Getenv("TEST_DB_URL"), "test") {
		t.Fatal("TEST_DB_URL does not appear to be a test database")
	}
	t.Cleanup(func() { db.Close() })
	return &apiConfig{db: database.New(db)}
}

func createTestUser(t *testing.T, cfg *apiConfig) database.User {
	t.Helper()
	hash := createTestHashPassword(t, "testpassword123")
	user, err := cfg.db.CreateUser(context.Background(), database.CreateUserParams{
		ID:           uuid.New(),
		Email:        fmt.Sprintf("%v@test.com", time.Now().UTC()),
		PasswordHash: hash,
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	})
	if err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}
	t.Cleanup(func() {
		cfg.db.DeleteUserByID(context.Background(), user.ID)
	})
	return user
}
func createTestHashPassword(t *testing.T, password string) string {
	t.Helper()
	hash, err := auth.HashPassword(password)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}
	return hash
}
func createTestWorkoutSession(t *testing.T, cfg *apiConfig, userID uuid.UUID) database.WorkoutSession {
	t.Helper()
	//TODO
	return database.WorkoutSession{}
}
