package main

import (
	"context"
	"errors"
	"testing"

	"github.com/DavidMWeaver4/Davids_Workout_App/internal/database"
	"github.com/google/uuid"
)

func TestAuthorizeWorkoutExercise_Success(t *testing.T) {
	x := newWorkoutExerciseFixture(t)

	gotWOExercise, err := x.cfg.authorizeWorkoutExercise(x.ctx, x.woExercise.ID, x.user.ID)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if gotWOExercise.ID != x.woExercise.ID {
		t.Fatalf("expected exercise ID %v, got %v", x.woExercise.ID, gotWOExercise.ID)
	}

	if gotWOExercise.WorkoutSessionID != x.session.ID {
		t.Fatalf("expected session ID %v, got %v", x.session.ID, gotWOExercise.WorkoutSessionID)
	}
}

func TestAuthorizeWorkoutExercise_Forbidden(t *testing.T) {
	x := newWorkoutExerciseFixture(t)

	otherUser := createTestUser(t, x.cfg)

	_, err := x.cfg.authorizeWorkoutExercise(x.ctx, x.woExercise.ID, otherUser.ID)

	if err == nil {
		t.Fatal("expected ErrForbidden, got nil")
	}

	if !errors.Is(err, ErrForbidden) {
		t.Fatalf("expected ErrForbidden, got %v", err)
	}
}

func TestAuthorizeWorkoutExercise_NotFound(t *testing.T) {
	x := newWorkoutExerciseFixture(t)

	_, err := x.cfg.authorizeWorkoutExercise(x.ctx, uuid.New(), x.user.ID)

	if err == nil {
		t.Fatal("expected ErrNotFound, got nil")
	}

	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

type woExerciseFixture struct {
	ctx        context.Context
	cfg        *apiConfig
	user       database.User
	session    database.WorkoutSession
	exercise   database.Exercise
	woExercise database.WorkoutExercise
}

func newWorkoutExerciseFixture(t *testing.T) woExerciseFixture {
	t.Helper()

	cfg := newTestAPIConfig(t)
	user := createTestUser(t, cfg)
	session := createTestWorkoutSession(t, cfg, user.ID)
	exercise := createTestExercise(t, cfg, user.ID)
	woExercise := createTestWorkoutExercise(t, cfg, session, exercise.ID)

	return woExerciseFixture{
		ctx:        context.Background(),
		cfg:        cfg,
		user:       user,
		session:    session,
		exercise:   exercise,
		woExercise: woExercise,
	}
}
