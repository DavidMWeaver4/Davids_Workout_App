package main

import (
	"context"
	"errors"
	"testing"

	"github.com/DavidMWeaver4/Davids_Workout_App/internal/database"
	"github.com/google/uuid"
)

func TestAuthorizeWorkoutSession_Success(t *testing.T) {
	x := newWorkoutSessionFixture(t)

	gotSession, err := x.cfg.authorizeWorkoutSession(x.ctx, x.session.ID, x.user.ID)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if gotSession.ID != x.session.ID {
		t.Errorf("expected session ID %v, got %v", x.session.ID, gotSession.ID)
	}

	if gotSession.UserID != x.user.ID {
		t.Errorf("expected user ID %v, got %v", x.user.ID, gotSession.UserID)
	}
}

func TestAuthorizeWorkoutSession_Forbidden(t *testing.T) {
	x := newWorkoutSessionFixture(t)
	otherUser := createTestUser(t, x.cfg)
	if x.user.ID == otherUser.ID {
		t.Fatal("test setup error: users have same ID")
	}
	_, err := x.cfg.authorizeWorkoutSession(x.ctx, x.session.ID, otherUser.ID)
	if err == nil {
		t.Fatal("expected ErrForbidden, got nil")
	}
	if !errors.Is(err, ErrForbidden) {
		t.Fatalf("expected ErrForbidden, got %v", err)
	}
}

func TestAuthorizeWorkoutSession_NotFound(t *testing.T) {
	x := newWorkoutSessionFixture(t)

	_, err := x.cfg.authorizeWorkoutSession(x.ctx, uuid.New(), x.user.ID)
	if err == nil {
		t.Fatal("expected ErrNotFound, got nil")
	}
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

type workoutSessionFixture struct {
	ctx     context.Context
	cfg     *apiConfig
	user    database.User
	session database.WorkoutSession
}

func newWorkoutSessionFixture(t *testing.T) workoutSessionFixture {
	t.Helper()

	cfg := newTestAPIConfig(t)
	user := createTestUser(t, cfg)
	session := createTestWorkoutSession(t, cfg, user.ID)

	return workoutSessionFixture{
		ctx:     context.Background(),
		cfg:     cfg,
		user:    user,
		session: session,
	}
}
