package main

import (
	"context"
	"errors"
	"testing"

	"github.com/DavidMWeaver4/Davids_Workout_App/internal/database"
	"github.com/google/uuid"
)

func TestAuthorizeCardioAndSets_Success(t *testing.T) {
	x := newCardioSetFixture(t)
	gotCardioSet, err := x.cfg.authorizeCardioSet(x.ctx, x.cardioSet.ID, x.user.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if gotCardioSet.ID != x.cardioSet.ID {
		t.Fatalf("expected cardio ID %v, got %v", x.cardioSet.ID, gotCardioSet.ID)
	}
	if gotCardioSet.WorkoutExercisesID != x.cardioSet.WorkoutExercisesID {
		t.Fatalf("expected exercise ID %v, got %v: ", x.cardioSet.WorkoutExercisesID, gotCardioSet.WorkoutExercisesID)
	}
}

func TestAuthorizedCardioAndSets_Forbidden(t *testing.T) {
	x := newCardioSetFixture(t)
	otherUser := createTestUser(t, x.cfg)
	_, err := x.cfg.authorizeCardioSet(x.ctx, x.cardioSet.ID, otherUser.ID)
	if err == nil {
		t.Fatal("expected ErrForbidden, got nil")
	}

	if !errors.Is(err, ErrForbidden) {
		t.Fatalf("expected ErrForbidden, got %v", err)
	}
}

func TestAuthorizedCardioAndSets_NotFound(t *testing.T) {
	x := newCardioSetFixture(t)
	_, err := x.cfg.authorizeCardioSet(x.ctx, uuid.New(), x.user.ID)
	if err == nil {
		t.Fatal("expected ErrNotFound, got nil")
	}
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

type cardioSetFixture struct {
	ctx        context.Context
	cfg        *apiConfig
	user       database.User
	session    database.WorkoutSession
	exercise   database.Exercise
	woExercise database.WorkoutExercise
	cardioSet  database.CardioAndSet
}

func newCardioSetFixture(t *testing.T) cardioSetFixture {
	t.Helper()
	cfg := newTestAPIConfig(t)
	user := createTestUser(t, cfg)
	session := createTestWorkoutSession(t, cfg, user.ID)
	exercise := createTestExercise(t, cfg, user.ID)
	woExercise := createTestWorkoutExercise(t, cfg, session, exercise.ID)
	cardioSet := createTestCardioSet(t, cfg, woExercise)
	return cardioSetFixture{
		ctx:        context.Background(),
		cfg:        cfg,
		user:       user,
		session:    session,
		exercise:   exercise,
		woExercise: woExercise,
		cardioSet:  cardioSet,
	}
}
