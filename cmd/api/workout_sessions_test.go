package main

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
)

func TestAuthorizeWorkoutSession_Success(t *testing.T) {
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

func TestAuthorizeWorkoutSession_Forbidden(t *testing.T) {
	ctx := context.Background()
	cfg := newTestAPIConfig(t)

	user := createTestUser(t, cfg)
	otherUser := createTestUser(t, cfg)
	if user.ID == otherUser.ID {
		t.Fatal("test setup error: users have same ID")
	}

	session := createTestWorkoutSession(t, cfg, user.ID)
	_, err := cfg.authorizeWorkoutSession(ctx, session.ID, otherUser.ID)
	if err == nil {
		t.Fatal("expected ErrForbidden, got nil")
	}
	if !errors.Is(err, ErrForbidden) {
		t.Fatalf("expected ErrForbidden, got %v", err)
	}
}

func TestAuthorizeWorkoutSession_NotFound(t *testing.T) {
	ctx := context.Background()
	cfg := newTestAPIConfig(t)
	user := createTestUser(t, cfg)
	_ = createTestWorkoutSession(t, cfg, user.ID)
	//want to make sure there is something inside the db even if its going to fail

	_, err := cfg.authorizeWorkoutSession(ctx, uuid.New(), user.ID)
	if err == nil {
		t.Fatal("expected ErrForbidden, got nil")
	}
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}
