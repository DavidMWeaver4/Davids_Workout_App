package main

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
)

func TestAuthorizeWorkoutExercise_Success(t *testing.T) {
	ctx := context.Background()
	cfg := newTestAPIConfig(t)
	user := createTestUser(t, cfg)
	session := createTestWorkoutSession(t, cfg, user.ID)
	exercise := createTestExercise(t, cfg, user.ID)
	woExercise := createTestWorkoutExercise(t, cfg, session, exercise.ID)

	gotWOExercise, err := cfg.authorizeWorkoutExercise(ctx, woExercise.ID, user.ID)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if gotWOExercise.ID != woExercise.ID {
		t.Fatalf("expected exercise ID %v, got %v", woExercise.ID, gotWOExercise.ID)
	}
	if gotWOExercise.WorkoutSessionID != session.ID {
		t.Fatalf("expected session ID %v, got %v", session.ID, gotWOExercise.WorkoutSessionID)
	}
}

func TestAuthorizeWorkoutExercise_Forbidden(t *testing.T) {
	ctx := context.Background()
	cfg := newTestAPIConfig(t)

	user := createTestUser(t, cfg)
	otherUser := createTestUser(t, cfg)

	session := createTestWorkoutSession(t, cfg, user.ID)
	exercise := createTestExercise(t, cfg, user.ID)
	woExercise := createTestWorkoutExercise(t, cfg, session, exercise.ID)

	_, err := cfg.authorizeWorkoutExercise(ctx, woExercise.ID, otherUser.ID)

	if err == nil {
		t.Fatal("expected ErrForbidden, got nil")
	}

	if !errors.Is(err, ErrForbidden) {
		t.Fatalf("expected ErrForbidden, got %v", err)
	}
}

func TestAuthorizeWorkoutExercise_NotFound(t *testing.T) {
	ctx := context.Background()
	cfg := newTestAPIConfig(t)

	user := createTestUser(t, cfg)

	_, err := cfg.authorizeWorkoutExercise(ctx, uuid.New(), user.ID)

	if err == nil {
		t.Fatal("expected ErrNotFound, got nil")
	}

	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}
