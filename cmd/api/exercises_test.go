package main

import (
	"context"
	"errors"
	"testing"

	"github.com/DavidMWeaver4/Davids_Workout_App/internal/database"
	"github.com/google/uuid"
)

func TestAuthorizeExercise_Success(t *testing.T) {
	x := newExerciseFixture(t)
	gotExercise, err := x.cfg.authorizeExercise(x.ctx, x.exercise.ID, x.user.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if gotExercise.ID != x.exercise.ID {
		t.Fatalf("expected exercise ID %v, got %v", x.exercise.ID, gotExercise.ID)
	}
	if !gotExercise.UserID.Valid {
		t.Fatalf("expected exercise to have a valid UserID, got null")
	}
	if gotExercise.UserID.UUID != x.user.ID {
		t.Fatalf("expected exercise UserID %v, got %v", x.user.ID, gotExercise.UserID.UUID)
	}
}

func TestAuthorizeExercise_Forbidden(t *testing.T) {
	x := newExerciseFixture(t)
	otherUser := createTestUser(t, x.cfg)
	_, err := x.cfg.authorizeExercise(x.ctx, x.exercise.ID, otherUser.ID)
	if err == nil {
		t.Fatal("expected ErrForbidden, got nil")
	}
	if !errors.Is(err, ErrForbidden) {
		t.Fatalf("expected ErrForbidden, got %v", err)
	}
}

func TestAuthorizeExercise_NotFound(t *testing.T) {
	x := newExerciseFixture(t)
	_, err := x.cfg.authorizeExercise(x.ctx, uuid.New(), x.user.ID)
	if err == nil {
		t.Fatal("expected ErrNotFound, got nil")
	}
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

type exerciseFixture struct {
	ctx      context.Context
	cfg      *apiConfig
	user     database.User
	exercise database.Exercise
}

func newExerciseFixture(t *testing.T) exerciseFixture {
	t.Helper()
	cfg := newTestAPIConfig(t)
	user := createTestUser(t, cfg)
	exercise := createTestExercise(t, cfg, user.ID)
	return exerciseFixture{
		ctx:      context.Background(),
		cfg:      cfg,
		user:     user,
		exercise: exercise,
	}
}
